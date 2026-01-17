from time import sleep
from typing import List, Tuple
import polyline
import requests
from fastapi import FastAPI
import threading


ELEFTHERIA_COORDS = (35.17022784728593, 33.35889554051766)
S1_END_COORDS = (35.13160429484031, 33.299296813161504)
S2_END_COORDS = (35.11338633948102, 33.33255319068168)


def get_route(start, end):
    url = (
        "http://router.project-osrm.org/route/v1/driving/"
        f"{start[1]},{start[0]};{end[1]},{end[0]}"
        "?overview=full&geometries=polyline"
    )
    r = requests.get(url)
    r.raise_for_status()
    data = r.json()
    return polyline.decode(data["routes"][0]["geometry"])

def get_route_from_stops(stops):
    route = []
    for i in range(len(stops) - 1):
        route.extend(get_route(stops[i], stops[i + 1]))
    return route

def bus_update_worker(interval_seconds=2):
    while True:
        for bus in buses_in_service:
            bus.move_next()
        sleep(interval_seconds)


class Bus:
    def __init__(self, name: str, stops: List[Tuple[float, float]]):
        self.name = name
        self.stops = stops
        self.route = get_route_from_stops(stops)
        self.pos_index = 0          # <- add this
        self.pos = self.route[0]    # initial position

    def to_dict(self):
        return {
            "name": self.name,
            "stops": self.stops,
            "route": self.route,
        }

    def add_stop(self, new_stop: Tuple[float, float]):
        def dist(a, b):
            return (a[0] - b[0]) ** 2 + (a[1] - b[1]) ** 2

        min_increase = float("inf")
        insert_index = len(self.stops)

        for i in range(len(self.stops) - 1):
            increase = (
                dist(self.stops[i], new_stop)
                + dist(new_stop, self.stops[i + 1])
                - dist(self.stops[i], self.stops[i + 1])
            )
            if increase < min_increase:
                min_increase = increase
                insert_index = i + 1

        self.stops.insert(insert_index, new_stop)
        self.route = get_route_from_stops(self.stops)

    def move_next(self):
        if not self.route:
            return  # safety check
        self.pos_index += 1
        if self.pos_index >= len(self.route):
            self.pos_index = 0  # loop back to start
        self.pos = self.route[self.pos_index]


buses_in_service : List[Bus] = []

def find_closest_bus(origin: Tuple[float, float], buses: List[Bus]) -> Bus:
    def dist(a, b):
        return (a[0] - b[0])**2 + (a[1] - b[1])**2  # squared distance is enough

    closest_bus = None
    min_distance = float("inf")

    for bus in buses:
        # check all points along the bus route
        for point in bus.route:
            d = dist(origin, point)
            if d < min_distance:
                min_distance = d
                closest_bus = bus

    return closest_bus


# API endpoints
app = FastAPI()
@app.get("/get_buses") # for dashboard
def get_buses():
    return [bus.to_dict() for bus in buses_in_service]

@app.get("/request_transfer")
def request_transfer(origin_str:str, dest_str: str):
    olat, olon = map(float, origin_str.split(","))
    origin = (olat, olon)
    dlat, dlon = map(float, dest_str.split(","))
    dest = (dlat, dlon)

    # locate the closest bus to origin
    closest_bus = find_closest_bus(origin, buses_in_service)
    closest_bus.add_stop(origin)
    closest_bus.add_stop(dest)
    return closest_bus.name


@app.on_event("startup")
def startup():

    # add all service areas
    # service area 1
    bus_s1 = Bus(
        name = "S1",
        stops = [ELEFTHERIA_COORDS,S1_END_COORDS]
    )

    buses_in_service.append(bus_s1)

    # service area 2
    bus_s2 = Bus(
        name = "S2",
        stops = [ELEFTHERIA_COORDS,S2_END_COORDS]
    )
    buses_in_service.append(bus_s2)

    # start bus positional updates
    threading.Thread(target=bus_update_worker, daemon=True).start()



import random
from time import sleep
from typing import List, Tuple
import polyline
import requests
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
import threading


ELEFTHERIA_COORDS = (35.17022784728593, 33.35889554051766)
S1_END_COORDS = (35.13160429484031, 33.299296813161504)
S2_END_COORDS = (35.11338633948102, 33.33255319068168)
S3_END_COORDS = (35.12984920055877, 33.36299761491542)
S4_END_COORDS = (35.14966690544324, 33.41059674652208)

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
        route.extend(get_route(stops[i][1:], stops[i + 1][1:]))
    return route

def bus_update_worker(interval_seconds=1.2):
    while True:
        for bus in buses_in_service:
            bus.move_next()
        sleep(interval_seconds)


class Bus:
    def __init__(self, name: str, stops: List[Tuple[str, float, float]]):
        self.name = name
        self.stops = stops
        self.route = get_route_from_stops(stops)
        self.pos_index = 0
        self.pos = self.route[0]
        self.next_stop_idx = 1  # The stop we are currently driving toward

    def add_stop(self, new_stop: Tuple[str, float, float]):
        def dist_sq(a, b):
            return (a[1] - b[1]) ** 2 + (a[2] - b[2]) ** 2

        # find best insertion point
        min_increase = float("inf")
        insert_index = len(self.stops)

        for i in range(self.next_stop_idx - 1, len(self.stops) - 1):
            increase = (dist_sq(self.stops[i], new_stop) +
                        dist_sq(new_stop, self.stops[i + 1]) -
                        dist_sq(self.stops[i], self.stops[i + 1]))
            if increase < min_increase:
                min_increase = increase
                insert_index = i + 1

        self.stops.insert(insert_index, new_stop)

        # recalculate route from current position
        leg_to_target = get_route(self.pos, self.stops[self.next_stop_idx][1:])
        remaining_legs = get_route_from_stops(self.stops[self.next_stop_idx:])

        self.route = leg_to_target + remaining_legs
        self.pos_index = 0  # Start at the beginning of this new path

    def move_next(self):
        if not self.route: return
        self.pos_index += 1

        if self.pos_index < len(self.route):
            self.pos = self.route[self.pos_index]

            # Check if we reached the target stop to increment index
            target = self.stops[self.next_stop_idx]
            if abs(self.pos[0] - target[1]) < 0.0001 and abs(self.pos[1] - target[2]) < 0.0001:
                if self.next_stop_idx < len(self.stops) - 1:
                    self.next_stop_idx += 1
        else:
            self.pos_index = 0  # Loop logic

    def to_dict(self):
        return {
            "name": self.name,
            "stops": [{"name": s[0], "lat": s[1], "lon": s[2]} for s in self.stops],
            "route": self.route,
            "pos": self.pos,
        }

buses_in_service: List[Bus] = []

def find_closest_bus(origin: Tuple[float, float], buses: List[Bus]) -> Bus:
    def dist(a, b):
        return (a[0] - b[0]) ** 2 + (a[1] - b[1]) ** 2

    closest_bus = None
    min_distance = float("inf")

    for bus in buses:
        for point in bus.route:
            d = dist(origin, point)
            if d < min_distance:
                min_distance = d
                closest_bus = bus

    return closest_bus


# API endpoints
app = FastAPI()
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


@app.get("/get_buses")
def get_buses():
    return [bus.to_dict() for bus in buses_in_service]


@app.get("/request_transfer")
def request_transfer(origin_str: str, dest_str: str):
    olat, olon = map(float, origin_str.split(","))
    dlat, dlon = map(float, dest_str.split(","))

    closest_bus = find_closest_bus((olat, olon), buses_in_service)
    stop_id = closest_bus.name + "-" + str(random.randint(10, 100))

    origin_stop = (stop_id+"(Pickup)", olat, olon)
    dest_stop = (stop_id, dlat, dlon)

    closest_bus.add_stop(origin_stop)
    closest_bus.add_stop(dest_stop)

    return stop_id


@app.on_event("startup")
def startup():
    bus_s1 = Bus(
        name="S1",
        stops=[
            ("Eleftheria", float(ELEFTHERIA_COORDS[0]), float(ELEFTHERIA_COORDS[1])),
            ("S1 End", float(S1_END_COORDS[0]), float(S1_END_COORDS[1])),
        ],
    )
    buses_in_service.append(bus_s1)

    bus_s2 = Bus(
        name="S2",
        stops=[
            ("Eleftheria", float(ELEFTHERIA_COORDS[0]), float(ELEFTHERIA_COORDS[1])),
            ("S2 End", float(S2_END_COORDS[0]), float(S2_END_COORDS[1])),
        ],
    )
    buses_in_service.append(bus_s2)

    bus_s3 = Bus(
        name="S3",
        stops=[
            ("Eleftheria", float(ELEFTHERIA_COORDS[0]), float(ELEFTHERIA_COORDS[1])),
            ("S3 End", float(S3_END_COORDS[0]), float(S3_END_COORDS[1])),
        ],
    )
    buses_in_service.append(bus_s3)

    bus_s4 = Bus(
        name="S4",
        stops=[
            ("Eleftheria", float(ELEFTHERIA_COORDS[0]), float(ELEFTHERIA_COORDS[1])),
            ("S4 End", float(S4_END_COORDS[0]), float(S4_END_COORDS[1])),
        ],
    )
    buses_in_service.append(bus_s4)

    threading.Thread(target=bus_update_worker, daemon=True).start()

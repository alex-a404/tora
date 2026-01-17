import streamlit as st
import folium
import pydeck as pdk
import requests
import math
from streamlit_folium import st_folium
from streamlit_autorefresh import st_autorefresh

# configs
GET_BUSES_ENDPOINT = "http://127.0.0.1:8000/get_buses"
REQ_ENDPOINT = "http://127.0.0.1:8000/request_transfer"
CENTER_COORDS = [35.16998609756835, 33.3608189662158]
AVG_SPEED_KMH = 30  # Assumed average speed for PoC
ARRIVAL_THRESHOLD_KM = 0.1  # 100 meters to consider "arrived" at stop

COLORS_MAP = {
    "S1": [255, 0, 0, 200],
    "S2": [0, 255, 0, 200],
    "S3": [0, 0, 255, 200]
}

st.set_page_config(layout="wide", page_title="Tora Bus Tracking")


for key, default in [("origin", None), ("dest", None), ("bus_requested", False), ("stop_id", None)]:
    if key not in st.session_state:
        st.session_state[key] = default



def calculate_eta(coord1, coord2, speed_kmh=AVG_SPEED_KMH):
    if not coord1 or not coord2:
        return 0, 0
    lat1, lon1, lat2, lon2 = map(math.radians, [coord1[0], coord1[1], coord2[0], coord2[1]])
    dlat, dlon = lat2 - lat1, lon2 - lon1
    a = math.sin(dlat / 2) ** 2 + math.cos(lat1) * math.cos(lat2) * math.sin(dlon / 2) ** 2
    distance = 6371 * (2 * math.asin(math.sqrt(a)))
    return distance, math.ceil((distance / speed_kmh) * 60)


def get_assigned_bus_data(stop_id):
    try:
        assigned_name = stop_id.split("-")[0]
        r = requests.get(GET_BUSES_ENDPOINT, timeout=5)
        all_buses = r.json()
        return next((bus for bus in all_buses if bus["name"] == assigned_name), None)
    except:
        return None


if st.session_state.bus_requested:
    st_autorefresh(interval=3000, key="bus_tracking_refresh")

    st.title("Live Bus Tracking")

    bus_data = get_assigned_bus_data(st.session_state.stop_id)

    if bus_data:
        bus_pos = bus_data["pos"]
        dist_to_origin, eta_to_origin = calculate_eta(bus_pos, st.session_state.origin)
        dist_to_dest, eta_to_dest = calculate_eta(bus_pos, st.session_state.dest)

        if dist_to_origin > ARRIVAL_THRESHOLD_KM:
            status_msg = f"Bus **{bus_data['name']}** is on its way to pick you up."
            metric_label = "ETA to Pickup"
            current_eta = eta_to_origin
        else:
            status_msg = "You are on board! Heading to your destination."
            metric_label = "ETA to Destination"
            current_eta = eta_to_dest

        # UI Display
        col1, col2, col3 = st.columns([2, 1, 1])
        with col1:
            st.success(f"**Stop ID:** {st.session_state.stop_id}")
            st.write(status_msg)
        with col2:
            st.metric(label=metric_label, value=f"{current_eta} min")
        with col3:
            if st.button("Cancel / New Request"):
                st.session_state.bus_requested = False
                st.rerun()

        # pydeck
        icon_settings = {"url": "https://img.icons8.com/ios-filled/100/000000/bus.png", "width": 128, "height": 128,
                         "anchorY": 128}

        bus_layer = pdk.Layer(
            "IconLayer",
            data=[{"coordinates": [bus_pos[1], bus_pos[0]], "name": bus_data["name"], "icon_item": icon_settings}],
            get_icon="icon_item", get_size=45, get_position="coordinates", pickable=True
        )

        route_layer = pdk.Layer(
            "PathLayer",
            data=[{"path": [[p[1], p[0]] for p in bus_data["route"]],
                   "color": COLORS_MAP.get(bus_data["name"], [128, 128, 128])}],
            get_path="path", get_color="color", width_scale=20, width_min_pixels=3
        )

        points_layer = pdk.Layer(
            "ScatterplotLayer",
            data=[
                {"name": "Pickup", "coords": [st.session_state.origin[1], st.session_state.origin[0]],
                 "color": [0, 255, 0]},
                {"name": "Destination", "coords": [st.session_state.dest[1], st.session_state.dest[0]],
                 "color": [255, 0, 0]}
            ],
            get_position="coords", get_color="color", get_radius=40,
        )

        st.pydeck_chart(pdk.Deck(
            map_style=pdk.map_styles.CARTO_LIGHT,
            initial_view_state=pdk.ViewState(latitude=bus_pos[0], longitude=bus_pos[1], zoom=14),
            layers=[route_layer, points_layer, bus_layer],
            tooltip={"text": "{name}"}
        ))
    else:
        st.warning("Waiting for live bus coordinates...")
        if st.button("Reset"):
            st.session_state.bus_requested = False
            st.rerun()

else:
    # requestor ui
    st.title("Request Bus(Demo) - Tora for Nicosia")

    if not st.session_state.origin:
        st.info("Click on the map to set your **Origin**")
    elif not st.session_state.dest:
        st.info("Click on the map to set your **Destination**")
    else:
        st.warning("Locations set. Ready to confirm request!")

    m = folium.Map(location=CENTER_COORDS, zoom_start=13)
    if st.session_state.origin:
        folium.Marker(st.session_state.origin, popup="Origin", icon=folium.Icon(color="green")).add_to(m)
    if st.session_state.dest:
        folium.Marker(st.session_state.dest, popup="Destination", icon=folium.Icon(color="red")).add_to(m)

    map_data = st_folium(m, height=500, width="100%")
    clicked = map_data.get("last_clicked")

    if clicked:
        curr_coords = [clicked["lat"], clicked["lng"]]
        if not st.session_state.origin:
            st.session_state.origin = curr_coords
            st.rerun()
        elif not st.session_state.dest:
            st.session_state.dest = curr_coords
            st.rerun()

    if st.session_state.origin and st.session_state.dest:
        if st.button("Confirm Request", use_container_width=True):
            try:
                params = {
                    "origin_str": f"{st.session_state.origin[0]},{st.session_state.origin[1]}",
                    "dest_str": f"{st.session_state.dest[0]},{st.session_state.dest[1]}"
                }
                response = requests.get(REQ_ENDPOINT, params=params)
                response.raise_for_status()
                st.session_state.stop_id = response.text.strip('"')
                st.session_state.bus_requested = True
                st.rerun()
            except Exception as e:
                st.error(f"Server error: {e}")
import streamlit as st
import pydeck as pdk
import requests
from streamlit_autorefresh import st_autorefresh

# Auto-refresh every 3 seconds
st_autorefresh(interval=3000, key="bus_refresh")
st.set_page_config(layout="wide")
st.title("Tora for Nicosia (Live)")

CENTER_COORDS = [33.3608189662158, 35.16998609756835]  # [Lon, Lat]

# Fetch bus data
def fetch_buses():
    try:
        r = requests.get("http://127.0.0.1:8000/get_buses", timeout=5)
        r.raise_for_status()
        return r.json()
    except:
        return []

data = fetch_buses()

# Fixed colors per bus
colors_map = {
    "S1": [255, 0, 0, 200],
    "S2": [0, 255, 0, 200],
    "S3": [0, 0, 255, 200]
}

bus_positions = []
route_paths = []
stop_positions = []

for bus in data:
    # Bus position
    lat, lon = bus["pos"] if isinstance(bus["pos"], (list, tuple)) else map(float, bus["pos"].split(","))
    bus_positions.append({
        "name": bus["name"],
        "coordinates": [lon, lat],
        "color": colors_map.get(bus["name"], [128, 128, 128, 200]),
        "icon_data": {"url": "https://img.icons8.com/ios-filled/50/000000/bus.png",
                      "width": 128, "height": 128, "anchorY": 128}
    })

    # Bus route
    path_coords = [[p[1], p[0]] for p in bus["route"]]
    route_paths.append({
        "name": bus["name"],
        "path": path_coords,
        "color": colors_map.get(bus["name"], [128, 128, 128, 200])
    })

    # Stops
    for stop in bus["stops"]:
        stop_positions.append({
            "name": stop["name"],
            "coordinates": [stop["lon"], stop["lat"]],
            "icon_data": {"url": "https://img.icons8.com/ios-filled/50/000000/marker.png",
                          "width": 128, "height": 128, "anchorY": 128}
        })

# Layers
routes_layer = pdk.Layer("PathLayer", data=route_paths, get_path="path", get_color="color", width_scale=20, width_min_pixels=2)
buses_layer = pdk.Layer("IconLayer", data=bus_positions, get_icon="icon_data", get_size=10, get_position="coordinates")
stops_layer = pdk.Layer("IconLayer", data=stop_positions, get_icon="icon_data", get_size=10, get_position="coordinates")

# Render map
st.pydeck_chart(pdk.Deck(
    map_style=pdk.map_styles.CARTO_LIGHT,
    initial_view_state=pdk.ViewState(latitude=CENTER_COORDS[1], longitude=CENTER_COORDS[0], zoom=13),
    layers=[routes_layer, buses_layer, stops_layer],
    tooltip={"text": "{name}"}
))

# Legend
st.markdown("### Bus Legend")
for bus_name, color in colors_map.items():
    hex_color = '#%02x%02x%02x' % tuple(color[:3])
    st.markdown(f'<div style="display:flex;align-items:center;">'
                f'<div style="width:20px;height:20px;background:{hex_color};margin-right:6px;"></div>'
                f'{bus_name}</div>', unsafe_allow_html=True)

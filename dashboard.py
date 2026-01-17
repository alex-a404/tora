import folium
import requests
import streamlit as st
from streamlit_folium import st_folium
from streamlit_autorefresh import st_autorefresh

# refresh every 3 seconds
st_autorefresh(interval=3000, key="bus_refresh")

def fetch_buses(api_url="http://127.0.0.1:8000/get_buses"):
    r = requests.get(api_url, timeout=5)
    r.raise_for_status()
    return r.json()

CENTER_COORDS = (35.16998609756835, 33.3608189662158)

st.title("Tora for Nicosia")

buses = fetch_buses()

m = folium.Map(location=CENTER_COORDS, zoom_start=13)

for bus in buses:
    for lat, lon in bus["route"]:
        folium.CircleMarker(
            location=(lat, lon),
            radius=2,
        ).add_to(m)

st_folium(m, width=1200, height=600)

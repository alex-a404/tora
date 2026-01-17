import streamlit as st
import folium
import requests
from streamlit_folium import st_folium

REQ_ENDPOINT = "http://127.0.0.1:8000/request_transfer"

st.set_page_config(layout="wide")

# initial session state
if "origin" not in st.session_state:
    st.session_state.origin = None
if "dest" not in st.session_state:
    st.session_state.dest = None
if "bus_requested" not in st.session_state:
    st.session_state.bus_requested = False
if "stop_id" not in st.session_state:
    st.session_state.stop_id = None


if st.session_state.bus_requested:
    st.title("Success!")
    st.markdown(f"Bus Requested!")
    st.success(f"**Your Stop ID:** {st.session_state.stop_id}")

    if st.button("Request another ride"):
        st.session_state.origin = None
        st.session_state.dest = None
        st.session_state.bus_requested = False
        st.rerun()


else:
    st.title("Request Bus - Nicosia")

    # Instructions based on state
    if not st.session_state.origin:
        st.info("Select **Origin**")
    elif not st.session_state.dest:
        st.info("Select **Destination**")
    else:
        st.warning("Ready to request!")

    CENTER_COORDS = [35.16998609756835, 33.3608189662158]
    m = folium.Map(location=CENTER_COORDS, zoom_start=13)

    # Add Markers
    if st.session_state.origin:
        folium.Marker(st.session_state.origin, popup="Origin", icon=folium.Icon(color="green")).add_to(m)
    if st.session_state.dest:
        folium.Marker(st.session_state.dest, popup="Destination", icon=folium.Icon(color="red")).add_to(m)

    # Map click capture
    map_data = st_folium(m, height=500, width="100%")
    clicked = map_data.get("last_clicked")

    if clicked:
        curr_coords = [clicked["lat"], clicked["lng"]]
        if not st.session_state.origin:
            st.session_state.origin = curr_coords
            st.rerun()  # Rerun to update instructions and marker
        elif not st.session_state.dest:
            st.session_state.dest = curr_coords
            st.rerun()

    # Request Button
    if st.session_state.origin and st.session_state.dest:
        if st.button("Confirm Request", use_container_width=True):
            origin_str = f"{st.session_state.origin[0]},{st.session_state.origin[1]}"
            dest_str = f"{st.session_state.dest[0]},{st.session_state.dest[1]}"

            try:
                params = {"origin_str": origin_str, "dest_str": dest_str}
                response = requests.get(REQ_ENDPOINT, params=params)
                response.raise_for_status()

                # Update State
                st.session_state.stop_id = response.text.strip('"')  # Strip quotes if returned as string
                st.session_state.bus_requested = True
                st.rerun()  # Force rerun to clear the map view immediately
            except Exception as e:
                st.error(f"Failed to connect to server: {e}")
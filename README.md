# Tora
Demand-responsive transit services project (Cyprus Transport Hackathon).

## Overview:
This project consists of a backend witten in Go with REST endpoints for the frontend (HTML/JS). It aims to reduce the amount of empty
buses driving around the city at night, and to direct them where there is demand.

## Architecture
The main backend service, written in Go, tracks `*Bus`es, its attributes, `Stop`, and user requests, exposing REST API endpoints to get buses and request a diversion(ride). Overall state of the system is maintained by a `Service`.

Once a request comes in, Tora finds the closest `Stop` to the origin and destination, and communicates with an optimization service written in Python via gRPC in order to find the best place in a `Bus` `Itinerary` to add it.

Upon receiving a gRPC response, a `TrackingSession` is launched which updates the end user with the current location of the requested bus and the Stops assigned to the request.

## Examples


<img width="603" height="1311" alt="IMG_6927" src="https://github.com/user-attachments/assets/41a3b70a-c516-4390-80d6-e78cee6f2c4b" />
<img width="603" height="1311" alt="IMG_6928" src="https://github.com/user-attachments/assets/c82cf869-daac-4043-a3ce-caeccbe97a7e" />
<img width="1166" height="905" alt="Screenshot from 2026-02-23 23-47-53" src="https://github.com/user-attachments/assets/a61ccbb1-5a49-4343-a2fd-a4c55cbbcfb8" />

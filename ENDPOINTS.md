# HTTP endpoints

All paths are relative to `/api/v1`.

## System

| Method | Path | Access |
| --- | --- | --- |
| GET | /health | Public |

## Auth

| Method | Path | Access |
| --- | --- | --- |
| POST | /auth/register | Public |
| POST | /auth/login | Public |
| POST | /auth/refresh | Public |
| POST | /auth/logout | Public |
| GET | /auth/me | Authenticated |

## Warehouses

| Method | Path | Access |
| --- | --- | --- |
| GET | /warehouses | Public |
| GET | /warehouses/{id} | Public |
| POST | /warehouses | Admin |
| PATCH | /warehouses/{id} | Admin |
| PATCH | /warehouses/{id}/activate | Admin |
| DELETE | /warehouses/{id} | Admin |
| GET | /storage-zones | Worker, Logist, Admin |
| POST | /storage-zones | Admin |
| PATCH | /storage-zones/{id} | Admin |
| PATCH | /storage-zones/{id}/activate | Admin |
| DELETE | /storage-zones/{id} | Admin |
| GET | /gates | Worker, Logist, Admin |
| POST | /gates | Admin |
| PATCH | /gates/{id} | Admin |
| PATCH | /gates/{id}/activate | Admin |
| DELETE | /gates/{id} | Admin |
| GET | /product-types | Public |
| GET | /cargo-place-types | Public |

## Pickup calendar

| Method | Path | Access |
| --- | --- | --- |
| GET | /pickup-calendar | Client, Logist, Admin |
| POST | /pickup-calendar/blocks | Logist, Admin |
| DELETE | /pickup-calendar/blocks/{id} | Logist, Admin |
| PATCH | /pickup-calendar/capacity | Logist, Admin |

## Orders

| Method | Path | Access |
| --- | --- | --- |
| POST | /orders | Client |
| GET | /orders | Client, Logist, Admin |
| GET | /orders/{id} | Client, Logist, Admin |
| GET | /orders/{id}/history | Client, Logist, Admin |
| PATCH | /orders/{id}/cancel | Client, Admin |
| PATCH | /orders/{id}/status | Logist, Admin |

## Cargo items

| Method | Path | Access |
| --- | --- | --- |
| POST | /orders/{id}/cargo-items | Worker, Admin |
| GET | /cargo-items | Client, Worker, Logist, Admin |
| GET | /cargo-items/scan | Worker, Logist, Admin |
| GET | /cargo-items/{id} | Client, Worker, Logist, Admin |
| GET | /cargo-items/{id}/label | Worker, Logist, Admin |
| GET | /cargo-items/{id}/history | Client, Worker, Logist, Admin |
| PATCH | /cargo-items/{id}/status | Worker, Logist, Admin |
| PATCH | /cargo-items/{id}/assign-zone | Logist, Admin |
| PATCH | /cargo-items/{id}/assign-gate | Logist, Admin |

## Shipments

| Method | Path | Access |
| --- | --- | --- |
| POST | /shipments | Logist, Admin |
| GET | /shipments | Logist, Admin |
| GET | /shipments/{id} | Logist, Admin |
| POST | /shipments/{id}/items | Logist, Admin |
| DELETE | /shipments/{id}/items/{cargo_item_id} | Logist, Admin |
| PATCH | /shipments/{id}/status | Logist, Admin |

## Users

| Method | Path | Access |
| --- | --- | --- |
| GET | /users | Admin |
| POST | /users | Admin |
| PATCH | /users/{id} | Admin |
| PATCH | /users/{id}/block | Admin |
| DELETE | /users/{id} | Admin |

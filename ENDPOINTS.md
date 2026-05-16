# HTTP endpoints

All paths are relative to `/api/v1`.

## System

| Method | Path | Roles |
| --- | --- | --- |
| GET | /health | - |

## Auth

| Method | Path | Roles |
| --- | --- | --- |
| POST | /auth/register | - |
| POST | /auth/login | - |
| POST | /auth/refresh | - |
| POST | /auth/logout | - |
| GET | /auth/me | auth middleware |

## Warehouses

| Method | Path | Roles |
| --- | --- | --- |
| GET | /warehouses | - |
| GET | /warehouses/{id} | - |
| POST | /warehouses | Admin |
| PATCH | /warehouses/{id} | Admin |
| PATCH | /warehouses/{id}/activate | Admin |
| DELETE | /warehouses/{id} | Admin |
| GET | /storage-zones | - |
| POST | /storage-zones | Admin |
| PATCH | /storage-zones/{id} | Admin |
| PATCH | /storage-zones/{id}/activate | Admin |
| DELETE | /storage-zones/{id} | Admin |
| GET | /gates | - |
| POST | /gates | Admin |
| PATCH | /gates/{id} | Admin |
| PATCH | /gates/{id}/activate | Admin |
| DELETE | /gates/{id} | Admin |
| GET | /product-types | - |
| GET | /cargo-place-types | - |

## Pickup calendar

| Method | Path | Roles |
| --- | --- | --- |
| GET | /pickup-calendar | Client, Logist, Admin |
| POST | /pickup-calendar/blocks | Logist, Admin |
| DELETE | /pickup-calendar/blocks/{id} | Logist, Admin |
| PATCH | /pickup-calendar/capacity | Logist, Admin |

## Orders

| Method | Path | Roles |
| --- | --- | --- |
| POST | /orders | Client, Admin |
| GET | /orders | Client, Logist, Admin |
| GET | /orders/{id} | Client, Logist, Admin |
| GET | /orders/{id}/history | Client, Logist, Admin |
| PATCH | /orders/{id}/cancel | Client, Admin |
| PATCH | /orders/{id}/status | Logist, Admin |

## Cargo items

| Method | Path | Roles |
| --- | --- | --- |
| POST | /orders/{id}/cargo-items | Worker, Admin |
| GET | /cargo-items | Client, Worker, Logist, Admin |
| GET | /cargo-items/scan | Client, Worker, Logist, Admin |
| GET | /cargo-items/{id} | Client, Worker, Logist, Admin |
| GET | /cargo-items/{id}/label | Client, Worker, Logist, Admin |
| GET | /cargo-items/{id}/history | Client, Worker, Logist, Admin |
| PATCH | /cargo-items/{id}/status | Worker, Logist, Admin |
| PATCH | /cargo-items/{id}/assign-zone | Logist, Admin |
| PATCH | /cargo-items/{id}/assign-gate | Logist, Admin |

## Shipments

| Method | Path | Roles |
| --- | --- | --- |
| POST | /shipments | Logist, Admin |
| GET | /shipments | Logist, Admin |
| GET | /shipments/{id} | Logist, Admin |
| POST | /shipments/{id}/items | Logist, Admin |
| DELETE | /shipments/{id}/items/{cargo_item_id} | Logist, Admin |
| PATCH | /shipments/{id}/status | Logist, Admin |

## Users

| Method | Path | Roles |
| --- | --- | --- |
| GET | /users | Admin |
| POST | /users | Admin |
| PATCH | /users/{id} | Admin |
| PATCH | /users/{id}/block | Admin |
| DELETE | /users/{id} | Admin |


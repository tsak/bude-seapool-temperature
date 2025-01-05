# Bude Seapool Temperature

Read temperatures from an [iMonnit](https://www.imonnit.com/) temperature sensor embedded in the pool
and generate an image for digital signage by the pool.

![Live temperature](https://spt.tsak.dev/temperature.png)

## Public API

This service exposes two endpoints:

### Latest reading

`GET /api/v1/temperature`

```bash
$ curl -s https://spt.tsak.dev/api/v1/temperature
```

```json
{
  "temperature": 13.4,
  "datetime": "2024-11-07T22:30:00Z"
}
```

### List of last measurements

`GET /api/v1/temperatures`

```bash
$ curl -s https://spt.tsak.dev/api/v1/temperatures
```

```json
[
  {
    "temperature": 13.4,
    "datetime": "2024-11-07T22:30:00Z"
  },
  {
    "temperature": 13.6,
    "datetime": "2024-11-07T22:19:58Z"
  },
  {
    "temperature": 13.5,
    "datetime": "2024-11-07T22:09:58Z"
  }
]
```

## Prerequisites

- Go 1.23
- [Air](https://github.com/air-verse/air)
- iMonnit API details

## Setup

Copy `.env.sample` to `.env` and fill in `MONNIT_SENSOR_ID`, `MONNIT_API_KEY_ID` and `MONNIT_API_SECRET_KEY`


## Development

Run `air` and connect to [localhost:3001](http://localhost:3001)

```bash
# Continuously build and reload
air
```

## Building

Create `bude-seapool-temperature` binary

```bash
go build
```
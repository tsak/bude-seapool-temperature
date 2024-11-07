# Bude Seapool Temperature

Read temperatures from an [iMonnit](https://www.imonnit.com/) temperature sensor embedded in the pool
and generate an image for digital signage by the pool.

![Live temperature](https://spt.tsak.dev/temperature.png)

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
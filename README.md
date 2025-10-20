# NextMetro

Find the next [MetroTransit](https://www.metrotransit.org/) bus or train at a given stop!

## Usage
You can access the service from your web browser or through CURL.

```zsh
curl https://ajkn.dev/metro/GREEN/EABK/EAST
# There's a eastbound Green Train coming to east bank station in 4 Min.
```

### Variables
``` curl https://ajkn.dev/metro/{NAME}/{STATION}/{DIRECTION}
- {name}        The name of the METRO. This can either be the direct ID, like 901, or for the LRT or BRT lines their name, like BLUE.
- {station}     The ID of the station, like EABK.
- {direction}   The direction, like NORTH, SOUTH, EAST or WEST. This depends on the line itself.
```

## Locally 
```zsh
git clone https://github.com/ajkneisl/nextmetro
cd nextmetro
go build
./nextmetro
```
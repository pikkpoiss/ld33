#!/usr/bin/env bash

GITROOT=`git rev-parse --show-toplevel`

cd $GITROOT

build_aesprite() {
  DEST=${2:-tmp}
  mkdir -p ${DEST}
  aseprite \
    --batch assets/${1}.ase \
    --save-as ${DEST}/${1}_00.png
}

build_aesprite numbered_squares
build_aesprite special_squares
build_aesprite highlight
build_aesprite human01
build_aesprite skeleton01
build_aesprite box01
build_aesprite ghost01
build_aesprite spikes01
build_aesprite bubble
build_aesprite mouse
build_aesprite tiles
build_aesprite tiles "assets/tiled"
cp assets/*.png tmp/

TexturePacker \
  --format json-array \
  --trim-sprite-names \
  --trim-mode None \
  --size-constraints POT \
  --disable-rotation \
  --data src/resources/spritesheet.json \
  --sheet src/resources/spritesheet.png \
  tmp

rm -rf tmp
cd -

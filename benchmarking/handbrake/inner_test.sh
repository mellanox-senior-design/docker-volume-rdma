mkdir -p /data
mv bbb_sunflower_native_60fps_normal.mp4 /data/bbb.mp4
cd /data
HandBrakeCLI -i bbb.mp4 -o bbb_comp.mp4 --stop-at duration:60 -- Normal 2>&1 > log1

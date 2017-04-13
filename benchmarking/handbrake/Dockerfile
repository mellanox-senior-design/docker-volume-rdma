FROM ubuntu

RUN apt update && apt install -y wget handbrake handbrake-cli python
RUN wget http://distribution.bbb3d.renderfarming.net/video/mp4/bbb_sunflower_native_60fps_normal.mp4

ADD test.sh test.sh
ADD inner_test.sh inner_test.sh
ADD bench.py bench.py

CMD python bench.py

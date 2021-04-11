### Install postgres

git clone https://github.com/postgres/postgres.git
# remove -O2 flags from configure
./configure --enable-debug
make
sudo make install
# add /usr/local/pgsql/bin/postgres to PATH

### Create database

mkdir pg_data
initdb -D pg_data
pg_ctl -D pg_data -l logfile start
createdb csi
psql csi

### Set up Wireshark

* https://www.wireshark.org/#download
* If the UI says "You don't have permission to capture on local interfaces", click "installing ChmodBPF"

### Set up pg_filedump

git clone https://git.postgresql.org/git/pg_filedump.git
# fix compile errors
make
sudo make install
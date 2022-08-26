https://stackoverflow.com/questions/60943536/how-to-use-ssl-for-a-backend-ec2-instance-without-a-domain
https://caddyserver.com/docs/install#debian-ubuntu-raspbian
https://nip.io/

---

history
    1  git
    2  go
    3  sudo apt install golang-go
    4  sudo apt-get update
    5  sudo apt install golang-go
    6  go
    7  ls -al
    8  exit
    9  pwd
   10  mkdir app
   11  cd app
   12  pwd
   13  cd ~
   14  ls -al
   15  exit
   16  cd app
   17  ./main
   18  go main
   19  ls -al
   20  exit
   21  cd app
   22  ./main
   23  ./main &
   24  nohup ./main &
   25  ps aux
   26  ps aux | grep ./main
   27  kill 3163
   28  ps aux | grep ./main
   29  ls -al
   30  vi .env
   31  nohup ./main &
   32  make
   33  sudo apt install make
   34  gs
   35  make
   36  ps aux | grep ./main
   37  kill 3453
   38  ps aux | grep ./main
   39  kill 3456
   40  ps aux | grep ./main
   41  ps aux
   42  nohup ./main &
   43  ps aux | grep ./main
   44  nohup ./main &
   45  ps aux | grep ./main
   46  nohup ./main &
   47  cd app
   48  ps aux | grep ./main
   49  nohup ./main &
   50  ps aux | grep ./main
   51  nohup ./main &
   52  ps aux | grep ./main
   53  sudo touch /etc/systemd/system/animapu-api.service
   54  sudo vi /etc/systemd/system/animapu-api.service
   55  sudo systemctl daemon-reload
   56  sudo systemctl restart animapu-api
   57  systemctl status animapu-api.service
   58  cd ..
   59  ./app/main
   60  sudo vi /etc/systemd/system/animapu-api.service
   61  sudo systemctl daemon-reload
   62  sudo systemctl restart animapu-api
   63  sudo systemctl status animapu-api
   64  sudo vi /etc/systemd/system/animapu-api.service
   65  sudo systemctl restart animapu-api
   66  sudo systemctl status animapu-api
   67  sudo vi /etc/systemd/system/animapu-api.service
   68  sudo systemctl restart animapu-api
   69  sudo systemctl status animapu-api
   70  sudo systemctl start animapu-api
   71  sudo systemctl daemon-reload
   72  sudo systemctl start animapu-api
   73  sudo systemctl status animapu-api
   74  sudo vi /etc/systemd/system/animapu-api.service
   75  sudo systemctl status animapu-api
   76  sudo vi /etc/systemd/system/animapu-api.service
   77  sudo systemctl daemon-reload
   78  sudo systemctl start animapu-api
   79  sudo systemctl status animapu-api
   80  cd app
   81  ls -al
   82  sudo vi .env
   83  sudo systemctl restart animapu-api
   84  sudo systemctl status animapu-api
   85  sudo apt install -y debian-keyring debian-archive-keyring apt-transport-https
   86  curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | sudo gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg
   87  curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | sudo tee /etc/apt/sources.list.d/caddy-stable.list
   88  sudo apt update
   89  sudo apt install caddy
   90  touch Caddyfile
   91  vi Caddyfile
   92  caddy start
   93  vi Caddyfile
   94  caddy start
   95  vi Caddyfile
   96  caddy start
   97  caddy
   98  caddy adapt
   99  caddy fmt
  100  caddy adapt
  101  caddy fmt
  102  vi Caddyfile
  103  caddy adapt
  104  caddy
  105  caddy stop
  106  caddy start --help
  107  caddy start -config Caddyfile
  108  caddy adapt --help
  109  caddy adapt -config=Caddyfile
  110  sudo vi /etc/caddy/Caddyfile
  111  caddy adapt
  112  caddy stop
  113  caddy start -config=Caddyfile
  114  sudo vi /etc/caddy/Caddyfile
  115  caddy stop
  116  caddy start -config=Caddyfile
  117  caddy status
  118  caddy hrlp
  119  caddy help
  120  caddy run --help
  121  caddy run -config=Caddyfile
  122  caddy stop
  123  vi Caddyfile
  124  sudo vi /etc/caddy/Caddyfile
  125  caddy start
  126  caddy
  127  caddy adapt
  128  caddy start
  129  curl localhost:2019/config/
  130  curl localhost:2020/config/
  131  caddy run
  132  sudo setcap CAP_NET_BIND_SERVICE=+eip $(which caddy)
  133  caddy run
  134  caddy start
  135  caddy run
  136  caddy stop
  137  caddy run
  138  ps aux
  139  kill 5449
  140  caddy run
  141  history
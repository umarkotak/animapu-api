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
  142  caddy start
  143  sudo vi /etc/caddy/Caddyfile
  144  sudo cat /etc/caddy/Caddyfile
  145  caddy stop
  146  caddy
  147  apt-get
  148  apt-get list
  149  apt-get remove
  150  sudo apt-get remove caddy
  151  sudo apt-get clean
  152  cd /
  153  df -h
  154  cd ~
  155  df -h
  156  pwd
  157  sudo apt install certbot python3-certbot-nginx
  158  cd /etc/nginx/sites-available/
  159  ls -al
  160  sudo apt install nginx
  161  sudo ufw app list
  162  sudo ufw allow 'Nginx HTTP'
  163  sudo ufw app list
  164  sudo ufw status
  165  sudo ufw allow 'Nginx HTTP'
  166  sudo ufw status
  167  ufw
  168  sudo ufw help
  169  sudo ufw enable
  170  sudo ufw status
  171  sudo ufw app list
  172  sudo ufw allow 'Nginx HTTP'
  173  sudo ufw allow 'Nginx Full'
  174  sudo ufw status
  175  sudo ufw allow 'Nginx HTTPS'
  176  sudo ufw allow 'OpenSSH'
  177  sudo ufw status
  178  systemctl status nginx
  179  systemctl start nginx
  180  sudo su
  181  systemctl status nginx
  182  sudo systemctl start nginx.service
  183  sudo systemctl status nginx.service
  184  sudo /etc/init.d/apache2 stop
  185  sudo su
  186  cd /etc/init.d/
  187  ls -al
  188  ps aux
  189  sudo lsof -i -P -n | grep LISTEN
  190  ps aux | grep caddy
  191  kill -9 5522
  192  sudo lsof -i -P -n | grep LISTEN
  193  sudo systemctl start nginx.service
  194  sudo systemctl status nginx.service
  195  curl -4 icanhazip.com
  196  sudo systemctl status nginx.service
  197  sudo lsof -i -P -n | grep LISTEN
  198  cd /etc/nginx/sites-enabled/
  199  ls -al
  200  vi default
  201  cd /var/www/
  202  ls -al
  203  cd html/
  204  ls -al
  205  index.nginx-debian.html
  206  vi /etc/nginx/sites-enabled/default
  207  cd /etc/nginx/sites-available/
  208  ls -al
  209  vi default
  210  sudo vi default
  211  sudo systemctl restart nginx.service
  212  sudo systemctl status nginx.service
  213  sudo ufw allow http
  214  sudo ufw allow https
  215  curl localhost
  216  curl localhost:80
  217  sudo apt install certbot python3-certbot-nginx
  218  cd /etc/nginx/sites-available/
  219  ls -al
  220  vi animapu.site
  221  touch animapu.site
  222  sudo touch animapu.site
  223  sudo vi animapu.site
  224  sudo nginx -t
  225  sudo systemctl reload nginx
  226  sudo systemctl status nginx
  227  sudo ufw status
  228  sudo certbot --nginx -d animapu.site -d www.animapu.site
  229  sudo certbot --nginx -v -d animapu.site -d www.animapu.site
  230  sudo vi /etc/nginx/sites-available/animapu.site
  231  cd /etc/nginx/sites-available/
  232  ls -al
  233  vi default
  234  cd ..
  235  ls -al
  236  sudo ln -s /etc/nginx/sites-available/animapu.site /etc/nginx/sites-enabled/
  237  sudo vi /etc/nginx/nginx.conf
  238  sudo nginx -t
  239  sudo systemctl restart nginx
  240  cd /etc/nginx/sites-available/
  241  ls -al
  242  vi animapu.site
  243  sudo vi animapu.site
  244  sudo systemctl restart nginx
  245  sudo systemctl reload nginx
  246  sudo systemctl status nginx
  247  sudo systemctl stop nginx
  248  sudo systemctl status nginx
  249  sudo systemctl start nginx
  250  sudo systemctl status nginx
  251  cd /etc/nginx/sites-enabled
  252  ls -al
  253  rm -rf default
  254  sudo rm -rf default
  255  sudo systemctl restart nginx
  256  sudo systemctl status nginx
  257  cd ..
  258  ls -al
  259  cd sites-available/
  260  cod
  261  ls
  262  cd an
  263  ls -al
  264  sudo vi animapu.site
  265  sudo systemctl restart nginx
  266  sudo systemctl reload nginx
  267  sudo certbot renew --dry-run
  268  sudo systemctl status certbot.timer
  269  sudo certbot --nginx -d example.com -d www.example.com
  270  sudo certbot --nginx -v -d animapu.site -d www.animapu.site
  271  cd ~
  272  sudo vi /etc/nginx/sites-available/animapu.site
  273  sudo systemctl reload ngingx
  274  sudo systemctl status ngingx
  275  sudo systemctl reload nginx
  276  sudo vi /etc/nginx/sites-available/animapu.site
  277  sudo nginx -t
  278  sudo vi /etc/nginx/sites-available/animapu.site
  279  sudo nginx -t
  280  sudo systemctl reload nginx
  281  sudo vi /etc/nginx/sites-available/animapu.site
  282  sudo systemctl restart nginx
  283  sudo systemctl status nginx
  284  curl animapu.site
  285  cd /etc/nginx/sites-available/
  286  ls -al
  287  cp animapu.site api.animapu.site
  288  sudo cp animapu.site api.animapu.site
  289  vi animapu.site
  290  sudo vi animapu.site
  291  sudo vi api.animapu.site
  292  history | grep ln
  293  sudo ln -s /etc/nginx/sites-available/api.animapu.site /etc/nginx/sites-enabled/
  294  sudo nginx -t
  295  sudo systemctl restart nginx
  296  sudo systemctl status nginx
  297  history | grep cert
  298  sudo certbot --nginx -v -d animapu.site -d www.animapu.site -d api.animapu.site
  299  sudo certbot --nginx -v -d animapu.site -d www.animapu.site
  300  ls -al
  301  cat animapu.site
  302  cat api.animapu.site
  303  sudo certbot --nginx -v -d animapu.site -d www.animapu.site -d api.animapu.site
  304  vi api.animapu.site
  305  sudo vi api.animapu.site
  306  sudo certbot --nginx -v -d animapu.site -d www.animapu.site -d api.animapu.site -d www.api.animapu.site
  307  sudo certbot --nginx -v -d animapu.site -d www.animapu.site -d api.animapu.site
  308  exit
  309  history | grep syst
  310  sudo systemctl restart animapu-api
  311  sudo systemctl stop animapu-api
  312  sudo systemctl statusp animapu-api
  313  sudo systemctl status animapu-api
  314  sudo systemctl start animapu-api
  315  sudo systemctl status animapu-api
  316  exit
  317  cd /etc/nginx/sites-available/
  318  ls -al
  319  cat animapu.site
  320  cd /etc/nginx/
  321  cd sites-available/
  322  ls -al
  323  cat default
  324  sudo vi default
  325  sudo vi /etc/nginx/sites-available/animapu.site
  326  sudo systemd status nginx
  327  sudo systemctl status nginx
  328  sudo systemctl reload nginx
  329  sudo systemctl status nginx
  330  sudo vi /etc/nginx/sites-available/animapu.site
  331  sudo nginx -t
  332  sudo vi /etc/nginx/sites-available/animapu.site
  333  sudo nginx -t
  334  sudo systemctl reload nginx
  335  sudo systemctl status nginx
  336  exit
  337  cat /etc/nginx/sites-available/animapu.site
  338  vi /etc/nginx/sites-available/animapu.site
  339  sudo vi /etc/nginx/sites-available/animapu.site
  340  history | grep cert
  341  sudo certbot --nginx -v -d animapu.site -d www.animapu.site
  342  sudo vi /etc/nginx/sites-available/animapu.site
  343  nginix -t
  344  nginx -t
  345  sudo nginx -t
  346  sudo vi /etc/nginx/sites-available/animapu.site
  347  rm -rf /etc/nginx/sites-available/api.animapu.site
  348  sudo rm -rf /etc/nginx/sites-available/api.animapu.site
  349  sudo nginx -t
  350  sudo rm -rf /etc/nginx/sites-enabled/api.animapu.site
  351  sudo nginx -t
  352  sudo systemctl restart nginx
  353  sudo systemctl status nginx
  354  history | grep sysmted
  355  history | grep sysmtemd
  356  history
  357  sudo systemctl restart animapu-api
  358  sudo systemctl restart nginx
  359  exit
  360  df -h
  361  cd /
  362  df -h
  363  cd ~
  364  ls -al
  365  cd app/
  366  ls -al
  367  cat nohup.out
  368  ls
  369  ls -al
  370  cat .env
  371  exit
  372  journalctl -f
  373  ls -al
  374  cd app/
  375  ls -al
  376  rm -rf nohup.out Caddyfile
  377  journalctl -u animapu-api.service
  378  touch app.log
  379  ls -al
  380  chmod 777 app.log
  381  ls -al
  382  exit
  383  cd app
  384  tail -f app.log
  385  exit
  386  cd app/
  387  tail -n 10000 app.log
  388  rm -rf app.log
  389  ls -al
  390  cd app
  391  ls -al
  392  df -h
  393  tail app.log
  394  ls -al
  395  sudo systemctl stop animapu-api
  396  sudo nohup ./main &
  397  ls -al
  398  tail -f nohup.out
  399  sudo tail -f nohup.out
  400  ps aux
  401  pas aux | grep main
  402  ps aux | grep main
  403  kill -9 214307 214308 214309
  404  ps aux | grep main
  405  kill -9 214309
  406  sudo kill -9 214309
  407  ps aux | grep main
  408  sudo systemctl status animapu-api
  409  exit
  410  cd app
  411  ls -al
  412  rm -rf nohup.out
  413  sudo nohup ./main &
  414  tail nohup.out
  415  sudo tail nohup.out
  416  exit
  417  sudo tail nohup.out
  418  cd app
  419  sudo tail nohup.out
  420  exit
  421  cd app
  422  ls -al
  423  ps aux | grep main
  424  kill -9 214599
  425  sudo kill -9 214599
  426  ps aux | grep main
  427  exit
  428  pwd
  429  exit
  430  cd app/
  431  ls -al
  432  sudo nohup .main &
  433  ps aux | grep main
  434  cd app
  435  sudo nohup .main &
  436  ls -al
  437  sudo nohup ./main &
  438  exit
  439  cd app
  440  pwd
  441  ls -al
  442  ./home/ubuntu/app/main
  443  workdir ~/app
  444  cd ~/app
  445  exit
  446  ps aux | grep main
  447  exit
  448  ps aux | grep main
  449  exit
  450  ps aux | grep main
  451  make run
  452  ls -al
  453  exit
  454  cd app
  455  make run_aws
  456  ps aux | grep main
  457  exit
  458  ps aux | grep main
  459  sudo nohup ./main &
  460  cd app
  461  sudo nohup ./main &
  462  ls -al
  463  rm -rf nohup.out
  464  ls -al
  465  exit
  466  cd ~/app && sudo nohup ./main &
  467  exit
  468  cd ~/app && sudo nohup ./main &
  469  exit
  470  cd ~/app && sudo nohup ./main &
  471  cd app
  472  rm -rf nohup.out
  473  exit
  474  cd ~/app && sudo nohup ./main &
  475  cd app/
  476  rm -rf nohup.out
  477  exit
  478  ls -al
  479  cd app
  480  ls -al
  481  exit
  482  cd ~/app && sudo nohup ./main &
  483  cd app
  484  ls -al
  485  rm -rf nohup.out
  486  exit
  487  sudo apt update
  488  sudo apt install postgresql postgresql-contrib
  489  sudo systemctl start postgresql.service
  490  sudo -i -u postgres
  491  exit
  492  sudo apt install postgresql postgresql-contrib
  493  service postgresql status
  494  service postgresql start
  495  sudo
  496  sudo su
  497  service postgresql status
  498  psql
  499  sudo -u postgres psql
  500  createuser --interactive --pwprompt
  501  sudo -u postgres psql
  502  sudo -u postgres createuser --interactive --pwprompt
  503  ls -al
  504  cd app
  505  ls -al
  506  git clone https://github.com/umarkotak/dexpense-api.git
  507  cd dexpense-api/
  508  sudo apt install git curl libssl-dev libreadline-dev zlib1g-dev autoconf bison build-essential libyaml-dev libreadline-dev libncurses5-dev libffi-dev libgdbm-dev
  509  curl -fsSL https://github.com/rbenv/rbenv-installer/raw/HEAD/bin/rbenv-installer | bash
  510  cd ~
  511  ls -al
  512  vi .bashrc
  513  exit
  514  irb
  515  rbenv -v
  516  cd app/
  517  ls
  518  cd dexpense-api/
  519  ls -al
  520  cat .ruby-version
  521  rbenv install 2.7.2
  522  irb
  523  exit
  524  irb
  525  cd app/dexpense-api/
  526  rbenv local 2.7.2
  527  ruby -v
  528  cat ~/.bashrc
  529  vi ~/.bashrc
  530  exit
  531  irb
  532  vi .bashrc
  533  exit
  534  rbenv global 2.7.2
  535  irb
  536  cd app/dexpense-api/
  537  ls -al
  538  sudo apt install libpq-dev
  539  gem install bundler
  540  bundle install --path vendor/bundle
  541  exit
  542  cd app/dexpense-api/
  543  touch .env
  544  vi .env
  545  bin/rails db:create
  546  bin/rails db:migrate
  547  clear
  548  bin/rails s -p 3000
  549  history
  550  cd /etc/nginx/
  551  ls -al
  552  cat nginx.conf
  553  ls -al
  554  cd sites-available/
  555  ls -al
  556  cat animapu.site
  557  cat default
  558  vi default
  559  history
  560  ls -al
  561  vi default
  562  cp default default.bak
  563  sudo cp default default.bak
  564  sudo systemctl restart nginx
  565  sudo systemctl status nginx
  566  ls -al
  567  vi default
  568  sudo vi default
  569  cat animapu.site
  570  cd ~/app
  571  ls -al
  572  cat .env
  573  ps aux
  574  ls -al
  575  cat app.log
  576  clear
  577  ls -al
  578  cd /etc/nginx/
  579  ls -al
  580  cd sites-enabled/
  581  ls -al
  582  cd ..
  583  cd sites-available/
  584  ls -al
  585  touch bukukaskita-api.animapu.site
  586  sudo touch bukukaskita-api.animapu.site
  587  ls -al
  588  cd ..
  589  cd sites-enabled/
  590  ls -al
  591  ln -s bukukaskita-api.animapu.site ../sites-available/bukukaskita-api.animapu.site
  592  ls -al
  593  history
  594  sudo ln -s /etc/nginx/sites-available/bukukaskita-api.animapu.site /etc/nginx/sites-enabled/
  595  ls -al
  596  cd ..
  597  cd sites-available/
  598  ls -al
  599  cp default default.bak.2
  600  sudo cp default default.bak.2
  601  cat animapu.site
  602  history
  603  sudo certbot --nginx -v -d bukukaskita-api.animapu.site
  604  ls -al
  605  sudo vi bukukaskita-api.animapu.site
  606  sudo systemctl restart nginx
  607  sudo systemctl status nginx
  608  cd ~/app/dexpense-api/
  609  bin/rails s -p 3000
  610  cd /etc/nginx/sites-available/
  611  ls -al
  612  vi default
  613  sudo vi default
  614  ls -al
  615  vi default
  616  sudo vi default
  617  sudo systemctl restart nginx
  618  sudo systemctl status nginx
  619  sudo vi default
  620  sudo systemctl stop nginx
  621  sudo systemctl start nginx
  622  sudo vi default
  623  sudo vi bukukaskita-api.animapu.site
  624  history | grep nginx
  625  sudo systemctl reload nginx
  626  sudo systemctl config nginx
  627  cd ..
  628  vi nginx.conf
  629  ls -al
  630  sudo vi nginx.conf
  631  sudo lsof -i -P -n | grep LISTEN
  632  curl localhost:6001/health
  633  curl localhost/health
  634  ls -al
  635  cd sites-available/
  636  ls -al
  637  cat bukukaskita-api.animapu.site
  638  cat animapu.site
  639  sudo vi bukukaskita-api.animapu.site
  640  sudo systemctl restart nginx
  641  sudo systemctl status nginx
  642  sudo vi bukukaskita-api.animapu.site
  643  cd ~
  644  ls -al
  645  cd app/dexpense-api/
  646  bin/rails s -p 3000
  647  gs
  648  git pull --rebase origin master
  649  bin/rails s -p 3000
  650  ls -al
  651  cd log/
  652  ls -al
  653  cd ..
  654  exit
  655  ls -al
  656  pg_restore --verbose --clean --no-acl --no-owner -h localhost -U admin -d dexpense_development latest.dump
  657  psql
  658  psql -u admin
  659  psql -U admin
  660  psql -U admin -p admin123
  661  sudo su postgres
  662  pg_restore --verbose --clean --no-acl --no-owner -h localhost -U admin -d dexpense_development latest.dump
  663  sudo su postgres
  664  ls -al
  665  cd app/
  666  ls
  667  cd dexpense-api/
  668  sudo nohup bin/rails s -p 3000
  669  ls -al
  670  rm nohup.out
  671  ls -al
  672  ps aux
  673  sudo nohup bin/rails s -p 3000
  674  tail -f nohup.out
  675  ls -al
  676  sudo ls nohup.out
  677  sudo tail -f nohup.out
  678  bin/rails s -p 3000
  679  nohup bin/rails s -p 3000
  680  sudo nohup bin/rails s -p 3000
  681  ls -al
  682  sudo tail nohup.out
  683  sudo
  684  sudo su
  685  bin/rails s -p 3000 &
  686  ps aux
  687  cd app/dexpense-api/
  688  bin/rails s -p 3000
  689  cat /home/ubuntu/app/dexpense-api/tmp/pids/server.pid
  690  kill -9 777200
  691  cat /home/ubuntu/app/dexpense-api/tmp/pids/server.pid
  692  bin/rails s -p 3000
  693  bg
  694  jobs
  695  exit
  696  cd ~/app && sudo nohup ./main &
  697  ls -al
  698  cd ~/app
  699  nohup ./main &
  700  ls -al
  701  rm -rf nohup.out
  702  nohup ./main &
  703  ls -al
  704  tail nohup.out
  705  ps aux
  706  kill -9 1063250
  707  kill -9 1063251
  708  ps aux
  709  nohup ./main &
  710  tail nohup.out
  711  tail -f nohup.out
  712  rm -rf nohup.out
  713  ls -al
  714  exit
  715  cd ~/app
  716  ps aux
  717  sudo nohup ./main &
  718  ls -al
  719  rm -rf nohup.out
  720  ls -al
  721  exit
  722  cd ~/app
  723  sudo nohup ./main &
  724  rm -rf nohup.out
  725  exit
  726  history
  727  cd /etc/nginx/
  728  ls -al
  729  sudo vi nginx.conf
  730  cd sites-available/
  731  ls -al
  732  sudo vi bukukaskita-api.animapu.site
  733  touch uler-tangga-api.animapu.site
  734  sudo touch uler-tangga-api.animapu.site
  735  ls -al
  736  sudo vi uler-tangga-api.animapu.site
  737  sudo certbot --nginx -v -d uler-tangga-api.animapu.site
  738  ls -al
  739  cd ..
  740  sudo vi nginx.conf
  741  exit
  742  ls -al
  743  cd app
  744  ls -al
  745  exit
  746  cd ~/app
  747  sudo nohup ./uler-tangga-api &
  748  rm -rf nohup.out
  749  ps aux
  750  ls -al
  751  sudo nohup ./main &
  752  rm -rf nohup.out
  753  ps aux
  754  history
  755  sudo systemctl reload nginx
  756  cd /etc/nginx/
  757  ls -al
  758  cd sites-enabled/
  759  ls -al
  760  history
  761  ln -s uler-tangga-api.animapu.site ../sites-available/uler-tangga-api.animapu.site
  762  ls -ak
  763  ls -al
  764  cd ..
  765  history
  766  history | grep ln
  767  sudo ln -s /etc/nginx/sites-available/uler-tangga-api.animapu.site /etc/nginx/sites-enabled/
  768  sudo systemctl restart nginx
  769  journalctl -xeu nginx.service
  770  ls -al
  771  cd sites-available/
  772  ls -al
  773  cat uler-tangga-api.animapu.site
  774  sudo vi uler-tangga-api.animapu.site
  775  sudo systemctl restart nginx
  776  exit
  777  cd ~/app
  778  sudo nohup ./uler-tangga-api &
  779  rm -rf nohup.out
  780  ps aux
  781  exit
  782  cd ~/app
  783  sudo nohup ./uler-tangga-api &
  784  rm -rf nohup.out
  785  ps aux
  786  exit
  787  nohup --help
  788  man nohup
  789  nohup --list
  790  nohup --help
  791  ps aux
  792  nohup
  793  nohup --help
  794  nohup --comand
  795  jobs
  796  jobs -l
  797  ps xw
  798  screen --help
  799  screen -list
  800  screen -r 777633.pts-2.ip-172-31-13-58
  801  exit
  802  cd app/
  803  ls -al
  804  cat app.log
  805  sudo nohup ./main &
  806  sudo pkill main
  807  ps aux
  808  sudo nohup ./main
  809  sudo nohup ./main &
  810  tail -f nohup.out
  811  sudo tail -f nohup.out
  812  rm -rf nohup.out
  813  ls -al
  814  cd ~/app
  815  sudo nohup ./main &
  816  rm -rf nohup.out
  817  exit
  818  cd ~/app
  819  sudo nohup ./main &
  820  rm -rf nohup.out
  821  ls -al
  822  exit
  823  cd app
  824  ls -al
  825  cat .env
  826  cd dexpense-api/
  827  ls -al
  828  cat .env
  829  psql
  830  sudo su postgres
  831  exit
  832  ls -al
  833  psql -d postgres
  834  pg_dump -U admin -d dexpense_development > dexpense.sql
  835  pg_dump -U postgres -d dexpense_development > dexpense.sql
  836  pg_dump -U admin -p admin -d dexpense_development > dexpense.sql
  837  pg_dump -U admin -a admin -d dexpense_development > dexpense.sql
  838  pg_dump --help
  839  pg_dump -U admin -W admin -d dexpense_development > dexpense.sql
  840  pg_dump -U admin -W -d dexpense_development > dexpense.sql
  841  cd app/
  842  ls -al
  843  cat .env
  844  cd dexpense-api/
  845  cat .env
  846  cd db/
  847  ls -al
  848  cd ..
  849  cd config/
  850  ls -al
  851  cat database.yml
  852  cd ..
  853  cat .env
  854  ls -al
  855  cat README.md
  856  cd ~
  857  ls -al
  858  rm -rf latest.dump
  859  rm -rf dexpense.sql
  860  pg_dump -U admin -W -d dexpense_development > dexpense.sql
  861  env
  862  sudo su
  863  clear
  864  pg_dump -d dexpense_development > dexpense.sql
  865  ls -al
  866  exit
  867  history
history -c
> ~/.bash_history
history -w

sudo truncate -s 0 /root/.bash_history

sudo truncate -s 0 /home/m1/.bash_history
sudo truncate -s 0 /home/kali/.bash_history
sudo truncate -s 0 /home/guest/.bash_history
sudo truncate -s 0 /home/agent0152/.bash_history
sudo truncate -s 0 /home/ivan/.bash_history

rm -f ~/.python_history ~/.mysql_history ~/.psql_history ~/.lesshst ~/.viminfo
sudo rm -f /root/.python_history /root/.viminfo

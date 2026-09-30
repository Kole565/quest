echo 'packer:pass' | sudo chpasswd # FIXME: comment out at prod

echo Creating user and readme
sudo /usr/sbin/useradd -m -s /bin/bash agent
echo 'agent:pass' | sudo chpasswd
echo 'Пароль это: test' | sudo tee /home/agent/readme

echo Configuring access rules
echo 'agent ALL=(ALL) NOPASSWD: /usr/bin/systemctl, /usr/bin/journalctl' | sudo tee /etc/sudoers.d/agent
sudo chmod 0440 /etc/sudoers.d/agent

echo Configuring local nginx server
sudo cp -r /tmp/static_server/* /var/www/html/
sudo chown -R www-data:www-data /var/www/html/
sudo systemctl stop nginx
sudo systemctl disable nginx

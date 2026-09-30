echo Configuring apt-proxy
echo 'Acquire::http::Proxy "http://10.0.2.2:3142";' | sudo tee /etc/apt/apt.conf.d/02apt-proxy
sudo sed -i 's|https://deb.debian.org|http://deb.debian.org|g' /etc/apt/mirrors/debian.list /etc/apt/mirrors/debian-security.list

echo Installing deps
sudo DEBIAN_FRONTEND=noninteractive apt-get update
sudo DEBIAN_FRONTEND=noninteractive apt-get install -y nano curl dnsutils nginx
sudo DEBIAN_FRONTEND=noninteractive apt-get install -y /tmp/config/jre8.deb

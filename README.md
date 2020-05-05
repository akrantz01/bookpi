# BookPi

Turn an old book into a covert server with a Raspberry Pi Zero and a battery for approximately $80 (USD).
In it's most basic form, it is just a Raspberry Pi Zero, a 6600 mAh LiIon battery, and a USB LiIon charger.
However, it can be extended to include whatever bells and whistles you want, like a display or a external HDD/SSD.

All the software is open source and can be found in this repo.
It is written in Python for the display, Golang for the server, and JavaScript for the web UI.

## Usage
Once you've setup your BookPi, connect to its Wi-Fi network and navigate to `book.pi` in the browser of your choice.
From there, you have access to a messaging system, file storage, and file sharing between users.

If you are looking to create your own, see the [Instructions](#instructions) section.

## Parts List

As stated above, there are infinitely many configurations you could build.
However, this is just the parts that I included used.

| Part | Cost | Quantity | Link |
|------|------|------|------|
| Raspberry Pi Zero W | $10 | 1 | [link](http://www.adafruit.com/product/3400) |
| 6600 mAh 3.7V Lithium Ion Battery Pack | $29.50 | 1 | [link](http://www.adafruit.com/product/353) |
| PowerBoost 1000 Charger | $19.95 | 1 | [link](http://www.adafruit.com/product/2465) |
| Monochrome 128x32 I2C OLED display | $17.50 | 1 | [link](http://www.adafruit.com/product/931) |
| Micro-B Male Plug USB Shell | $0.95 | 2 | [link](http://www.adafruit.com/product/1826) |
| Micro-B Female Plug USB Connector | $0.95 | 1 | [link](http://www.adafruit.com/product/1829) |
| Total Cost | $79.80

**Note**: The prices are as of May 2020

## Instructions
Follow the instructions below to build and setup your own BookPi.
I've broken the instructions into four parts:
- [Assembly](#assembly)
- [Setup](#setup)
- [Installation](#installation)
- [Configuration](#configuration)

### Assembly
To build your own, get yourself a hard or soft cover book that is at least 1 in or 2.5 cm thick.
The thicker book you have, the more space you'll have for components, but routing wires may be more difficult.

To start, I'd recommend laying out your parts on a sheet of paper roughly the size of your page.
Your diagram should include the Raspberry Pi Zero, PowerBoost 1000, LiPo battery, OLED Display, and charging port.
Be sure to consider where you'll be routing wires and the different heights of each component.
Once you're happy with your part layout, use an Exacto knife, or some other precision cutting tool to start slicing through your pages.
This is by far the most time-consuming and tedious part, so be patient.

Once all the components fit into your book, cut paths between each component for the wires.
Make sure you give yourself enough space for each component's wires and the various connectors.

### Setup
For wiring the display to the Raspberry Pi, see the [Adafruit guide](https://learn.adafruit.com/monochrome-oled-breakouts/python-wiring#adafruit-128x32-i2c-oled-display-11-6).
The header for the wiring diagram is: `Adafruit 128x32 I2C OLED Display`.
If you need a Raspberry Pi pinout, checkout the excelent [pinout.xyz](https://pinout.xyz).

To connect the charging port, Raspberry Pi, and PowerBoost 1000, you'll need to create some micro USB cables.
You'll need one cable going from female to male for charging the battery.
This cable will connect to the input micro USB terminal on the PowerBoost 1000.
You'll also need another male cable soldered onto the output of the PowerBoost 1000 that will connect to the Raspberry Pi.
These cables will only need to use the power pins, not the data or OTG pins.
Alternatively, you could cannibalize some micro USB cables if you want to save $2.85.

To connect the battery to the PowerBoost 1000, connect the JST connectors together and route the slack wherever.

Finally, install everything into your book.
This could be one of the trickier steps depending on how you laid out your parts.

### Installation
For the software, you have two options for getting it onto your Pi.
The easiest option is using the pre-made image, and the more involved would be downloading the installation bundle.
If you would like to compile this for yourself, see [BUILD.md](/BUILD.md).

#### The Pre-built Image (easier)
1. Download the `bookpi.img.gz` file from the [`Releases`](https://github.com/akrantz01/bookpi/releases/latest) tab
1. Decompress the image
    - use `gunzip` if you prefer the terminal
    - use the `7z` GUI if you prefer a nice interface
    - or use whatever program you're most comfortable with
1. Write the image onto the micro SD card with [Etcher](https://www.balena.io/etcher/)
1. Boot  up your BookPi and connect to it's Wi-Fi
    - Name/SSID: `BookPi`
    - Password: `changeme`

#### The Installation Bundle
1. Download the `bookpi.tar.gz` file from the [`Releases`](https://github.com/akrantz01/bookpi/releases/latest) tab
1. Copy it to your BookPi
1. Decompress the tarball with `tar zxf bookpi.tar.gz`
1. Change into the newly created `bookpi` directory
1. Run `sudo ./install.sh` to begin the installation process
    - ensure that you have an internet connection not using the built-in Wi-Fi card
1. Once the installation process is complete, the default password will be displayed and it will restart after 10 seconds.

### Configuration
To connect to your BookPi, connect to its Wi-Fi and SSH using into it with `book.pi` as the address, `pi` as the user, and `raspberry` as the password.
To configure your BookPi, there is a built in command `bookpi` to assist with the basic configuration.
If you need more advanced configuration, you'll need to manually edit the configuration files.

**Changing the Wi-Fi name/SSID**:<br/>
- `bookpi --ssid <SSID>`
  - Replace `<SSID>` with your new name

**Changing the Wi-Fi password**:<br/>
- `bookpi --password <PASSWORD>`
  - Replace `<PASSWORD>` with your new password
  - **Note**: the password must be at least 8 characters long

**Starting/stopping/restarting services**:<br/>
- Start: `bookpi --start`
- Stop: `bookpi --stop`
- Restart: `bookpi --restart`

**Configuration Files** (For advanced configuration only):
- `/etc/dhcpcd.conf` - for configuring [Dhcpcd](https://wiki.gentoo.org/wiki/Dhcpcd)
- `/etc/hostapd/hostapd.conf` - for configuring [Hostapd](https://wiki.gentoo.org/wiki/Hostapd)
- `/etc/dnsmasq.conf` - for configuring [Dnsmasq](https://wiki.gentoo.org/wiki/Dnsmasq)

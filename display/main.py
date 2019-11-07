import adafruit_ssd1306
from argparse import ArgumentParser
import board
import digitalio
from PIL import Image, ImageDraw, ImageFont
import statistics
import time

RUNNING = True
LINE_HEIGHT = 11

parser = ArgumentParser()
parser.add_argument("--width", action="store", default=128, type=int, help="Width of the display")
parser.add_argument("--height", action="store", default=32, type=int, help="Height of the display")
parser.add_argument("-m", "--mountpoint", action="store", default="/",
                    type=str, help="Mount point where data will be stored")
args = parser.parse_args()

# Connect to the display
oled = adafruit_ssd1306.SSD1306_I2C(args.width, args.height, board.I2C(), reset=digitalio.DigitalInOut(board.D4))

# Clear the oled
oled.fill(0)
oled.show()

# Image to store displayed text
image = Image.new('1', (args.width, args.height))

# Drawing context and font
draw = ImageDraw.Draw(image)
font = ImageFont.load_default()

# Loop for determining disk, network, and client statistics
while RUNNING:
    # Clear screen
    draw.rectangle((0, 0, args.width, args.height), outline=0, fill=0)

    # Get disk usage statistics
    disk = statistics.disk(args.mountpoint)

    # Draw disk usage text
    disk_text = f"Disk: {round(disk['used'], 1)}/{round(disk['total'], 1)} GB"
    draw.text((0, 0), disk_text, font=font, fill=255)

    # Draw IP address text
    ip_text = f"IP: {statistics.network()}"
    draw.text((0, LINE_HEIGHT), ip_text, font=font, fill=255)

    # Draw clients text
    clients_text = f"Clients: {statistics.clients()}"
    draw.text((0, LINE_HEIGHT*2), clients_text, font=font, fill=255)

    # Wait for 5 minutes before re-drawing
    time.sleep(5 * 60)

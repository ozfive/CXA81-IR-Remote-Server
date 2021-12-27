# CXA81-IR-Remote-Server

<!-- ABOUT THE PROJECT -->
## About The Project

I initially wanted to control my Cambridge Audio CXA81 amplifier from anywhere in the range of my home Wifi instead of the IR remote control. THen I figured this would be a great project to start open sourced from the beginning. I'm sure if you have come to this page you are looking for something for your CXA81 too. It's unfortunate that Cambridge Audio doesn't connect their StreamMagic app to their CXA series. Currently if you wanted to control your CXA81 you have to have a compatible CXN streamer. I don't have one. You probably don't either. ;)

There are a couple ways that you can connect to the CXA81 besides IR. I've found that IR is the quickest/easiest and most robust connection which includes volume control whereas 
the other methods such as RS232, and the Control Bus (Orange RCA plug on back panel)

I took some inspiration from [Turning a Cambridge Audio Amp On With A Phone](https://gaselli.software/2017/08/25/turning-a-cambridge-audio-amp-on-with-a-phone/) but 
ultimately had to brute force the codes with the IR Transmitter, a bash script that looped through [Claas Langbehn's RC-5 Template](http://lirc.sourceforge.net/remotes/rc-5/RC-5), a gigantic lirc config file, and some go old fashioned time. 


### Built With


The links here are just suggestions for the items used in this project. They do not contain any tracking or affiliate links as far as I know. I also don't claim that these are the least expensive products you could pick up either.

Hardware:
  * [Raspberry Pi Zero W](https://www.adafruit.com/product/3400)
  * [Gikfun Digital 38khz Ir Receiver Ir Transmitter Sensor Module Kit (Set of 3 each - 6 total)](https://www.amazon.com/dp/B0816P2545?psc=1&ref=ppx_yo2_dt_b_product_details)
  * [Jumper Wires - Connected 6" (F/F, 20 pack)](https://www.sparkfun.com/products/12796)


Software:

  * [Go Language](https://go.dev/)
  * [Gin Framework](https://github.com/gin-gonic/gin)

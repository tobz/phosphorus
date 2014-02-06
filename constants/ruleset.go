package constants

type ColorHandling uint8

const (
    ColorHandlingRvR  ColorHandling = 0x00
    ColorHandlingPvP  ColorHandling = 0x01
    ColorHandlingPvE  ColorHandling = 0x03
    ColorHandlingCoop ColorHandling = 0x04
)

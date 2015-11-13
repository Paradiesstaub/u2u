
# Go

## USB

* refresh USB list, if no device was initially found, when the user clicks on the drop-down widget.

* add UUID field to the device struct.


## ISO-Writer

* check before writing, if the usb device has enough capacity.

* check if the write target device is the same as listed in the drop-down by comparing the UUID of the drop-down item to the current UUID of the passed device-path.

* check what happens it the USB stick is removed during a write-procces.

#include "ZLED.h"
#include <iostream>
#include <gpiod.h>

ZLED::ZLED(int redPin, int greenPin, int bluePin)
    : redPin_(redPin), greenPin_(greenPin), bluePin_(bluePin),
      lastRed_(0), lastGreen_(0), lastBlue_(0),
      redLine_(nullptr), greenLine_(nullptr), blueLine_(nullptr),
      initSuccessful_(false) {

    struct gpiod_chip* chip = gpiod_chip_open_by_name("gpiochip0");
    if (!chip) {
        std::cerr << "ZLED Error: Failed to open gpiochip0. Continuing without hardware LEDs." << std::endl;
        return;
    }

    redLine_ = gpiod_chip_get_line(chip, redPin_);
    greenLine_ = gpiod_chip_get_line(chip, greenPin_);
    blueLine_ = gpiod_chip_get_line(chip, bluePin_);

    if (!redLine_ || !greenLine_ || !blueLine_) {
        std::cerr << "ZLED Error: Failed to get GPIO lines. Continuing without hardware LEDs." << std::endl;
        gpiod_chip_close(chip);
        return;
    }

    int flags = GPIOD_LINE_REQUEST_DIRECTION_OUTPUT;
    
    if (gpiod_line_request_output(redLine_, "ZScanner-LED", 0) < 0 ||
        gpiod_line_request_output(greenLine_, "ZScanner-LED", 0) < 0 ||
        gpiod_line_request_output(blueLine_, "ZScanner-LED", 0) < 0) {
        std::cerr << "ZLED Error: Failed to request GPIO lines as output." << std::endl;
        gpiod_chip_close(chip);
        return;
    }

    initSuccessful_ = true;
    
    Set(0, 0, 0);
}

ZLED::~ZLED() {
    if (initSuccessful_) {
        Set(0, 0, 0);
        if (redLine_) gpiod_line_release(redLine_);
        if (greenLine_) gpiod_line_release(greenLine_);
        if (blueLine_) gpiod_line_release(blueLine_);
    }
}

int ZLED::Set(uint32_t r, uint32_t g, uint32_t b) {
    if (initSuccessful_) {
        gpiod_line_set_value(redLine_, r > 0 ? 1 : 0);
        gpiod_line_set_value(greenLine_, g > 0 ? 1 : 0);
        gpiod_line_set_value(blueLine_, b > 0 ? 1 : 0);
    }
    
    lastRed_ = r;
    lastGreen_ = g;
    lastBlue_ = b;
    std::cout << "R:" << r << " G:" << g << " B:" << b << std::endl;
    return 0;
}

void ZLED::Get(uint32_t& r, uint32_t& g, uint32_t& b) const {
    r = lastRed_;
    g = lastGreen_;
    b = lastBlue_;
}

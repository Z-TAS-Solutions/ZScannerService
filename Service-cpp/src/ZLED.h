#pragma once

#include <cstdint>

struct gpiod_line;

class ZLED {
public:
    ZLED(int redPin, int greenPin, int bluePin);
    ~ZLED();

    int Set(uint32_t r, uint32_t g, uint32_t b);
    void Get(uint32_t& r, uint32_t& g, uint32_t& b) const;

private:
    int redPin_;
    int greenPin_;
    int bluePin_;
    uint32_t lastRed_;
    uint32_t lastGreen_;
    uint32_t lastBlue_;

    struct gpiod_line* redLine_;
    struct gpiod_line* greenLine_;
    struct gpiod_line* blueLine_;

    bool initSuccessful_;
};

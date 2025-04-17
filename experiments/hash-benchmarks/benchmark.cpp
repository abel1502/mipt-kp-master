#include <cstdio>
#include <cstdlib>
#include <memory>
#include <utility>
#include <string_view>

int main(int argc, const char **argv) {
    if (argc != 2) {
        printf(
            "Usage: %s <benchmark>\n"
            "\n"
            "Available benchmarks:\n"
            "  rabin\n"
            "  adler32\n"
            "  cyclic\n"
            "  md5\n"
            "  sha256\n"
            "  highway\n",
            argv[0]
        );
        return -1;
    }

    std::string_view benchmark = argv[1];

    if (benchmark == "rabin") {
        // TODO
        return 0;
    }

    if (benchmark == "adler32") {
        // TODO
        return 0;
    }

    if (benchmark == "cyclic") {
        // TODO
        return 0;
    }

    if (benchmark == "md5") {
        // TODO
        return 0;
    }

    if (benchmark == "sha256") {
        // TODO
        return 0;
    }

    if (benchmark == "highway") {
        // TODO
        return 0;
    }

    printf("Unknown benchmark: %s\n", argv[1]);
    return -1;
}

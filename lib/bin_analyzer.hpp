// Lib for pe files parsing

#ifndef BIN_ANALYZER_HPP
#define BIN_ANALYZER_HPP

#include <string>
#include <vector>

// Structure to represent a section in the binary file
struct Section {
    std::string name;
    std::vector<unsigned char> data;
};

// Class to analyze a binary file
class BinaryAnalyzer {
public:
    explicit BinaryAnalyzer(const std::string& filePath);
    void parseHeaders();
    Section getSection(const std::string& name); // Return a section by its name
    void writeSection(const Section& section);
    void save(const std::string& outputPath);

private:
    std::string filePath; // Path of the binary file
    // Internal data structures to hold parsed sections
    std::vector<Section> sections;
};

#endif // BIN_ANALYZER_HPP

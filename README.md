# ABCP Test Task

## Overview
This service is a Go application designed to simulate the generation and processing of tasks in a multi-threaded environment.

## Features
- **Asynchronous Task Generation and Processing**: Tasks are generated and processed concurrently using Go's goroutines.
- **Periodic Logging**: Logs the results of processed tasks every 3 seconds, separating successful and failed tasks.
- **Graceful Shutdown**: The application can be gracefully shut down using the context and the wait group, ensuring all goroutines are properly terminated.

## Getting Started

### Prerequisites

- Go (1.18 or higher)

### Installation

1. Clone the repository:
   ```sh
   git clone https://github.com/CroccoRush/ABCPTestTask.git
   cd ABCPTestTask
   ```

2. Build and run the application:
   ```sh
   go build
   ./ABCPTestTask
   ```

### Usage

The application will run for about 10 seconds, during which it will generate and process tasks. The results of the processed tasks will be logged every 3 seconds.

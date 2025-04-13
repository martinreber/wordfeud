# Letter Visualizer

## TODO's

- fix filter words in played words

This project is a simple web application that visualizes letter data from a Go application. It fetches data about letters, their counts, and values, and displays it in a user-friendly format.

## Project Structure

```
letter-visualizer
├── public
│   ├── index.html       # Main HTML document
│   ├── style.css        # Styles for the web application
│   └── script.js        # JavaScript code for fetching and displaying data
├── README.md            # Documentation for the project
└── package.json         # npm configuration file
```

## Setup Instructions

1. **Clone the repository:**
   ```
   git clone <repository-url>
   cd letter-visualizer
   ```

2. **Install dependencies:**
   ```
   npm install
   ```

3. **Run the Go application:**
   Make sure the Go application is running and accessible at `http://localhost:8080`.

4. **Start the web application:**
   You can use a simple HTTP server to serve the `public` directory. For example, you can use `http-server`:
   ```
   npx http-server public
   ```

5. **Open your browser:**
   Navigate to `http://localhost:8080` to view the application.

## Usage

- The application will display the current letters, their counts, and values.
- You can interact with the application to process input strings and see the updated letter data.

## Contributing

Feel free to submit issues or pull requests if you have suggestions or improvements for the project.
# etalert-backend

This is the installation guide for ETAlert-Backend project

## Requirements
- Go 1.18+
- Docker (optional, for containerization)

## Installation

1. Clone this repository:
    ```bash
    git clone https://github.com/Erxical/etalert-backend.git
    cd etalert-backend
    ```

2. Install dependencies:
    ```bash
    go mod tidy
    ```

3. Create an environment file:
    ```env
    # Google Map
    G_MAP_API_KEY=<PUT API_KEY HERE>

    # Gemini
    GEMINI_API_KEY=<PUT API_KEY HERE>

    # Azure Map
    AZURE_MAP_API_KEY=<PUT API_KEY HERE>

    # MongoDB
    MONGODB_URI=<PUT API_KEY HERE>

    # JWT
    JWT_SECRET=<SECRET>
    ```

   ### 3.1. Get Google Map API Key
   - Go to [Google Maps API](https://developers.google.com/maps)
   - Sign in to your account
   - Go to the **Keys & Credentials** tab on the left sidebar
   - Click on the **Create Credentials** button at the top left
   - Copy the API Key

   ### 3.2. Get Gemini API Key
   - Go to [Gemini API](https://aistudio.google.com/)
   - Sign in to your account
   - Click on the **Get API key** button at the top left of the sidebar
   - Click on **Create API key** button
   - Choose the project
   - Copy the API Key
  
   ### 3.3. Get Azure Map API Key
     - Go to [Azure Portal](https://portal.azure.com/#home)
     - Sign in to your account
     - Type 'Azure Maps Accounts' in the search bar and click on it
     - Click on the **Create** button
     - Choose Subscription, Resource group, name, and region
     - Check the agreement box at the bottom of the Basics tab
     - Click **Review + Create** then **Create**
     - Click on the account you just created
     - On the sidebar, click **Settings** then **Authentication**
     - Copy the Primary Key

    ### 3.4. Get MongoDB URI
     - Go to [MongoDB Atlas](https://www.mongodb.com/products/platform/atlas-database)
     - Click **Sign in**, then sign in with your account of choice
     - On the **Overview** tab, click **Create cluster**
     - Use the **Free Plan**
     - On the **Overview** tab under your cluster's name, click **Connect**
     - Choose **Drivers**
     - Set **Driver** to 'Go' and **Driver version** to '1.6 or later'
     - Copy the connection string

    ### 3.5. JWT Secret
     - JWT Secret can be any string

4. Run the application:
    ```bash
    go run main.go
    ```

5. To build a Docker image:
    ```bash
    docker build -t etalert-backend ./Dockerfile .
    ```

## Directory Structure
- `handler/`: Contains request handlers
- `middlewares/`: Middleware components for request processing
- `repository/`: Data storage layer and database interactions
- `service/`: Business logic of the application
- `validators/`: Input validation logic
- `websocket/`: WebSocket-related implementations


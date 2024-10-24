# NFTVC-profiles
NFTVC-Profiles is a service responsible for managing digital profiles within the NFT and Verifiable Credentials ecosystem. It supports operations for creating, updating, and retrieving profiles. Authorization and role management are handled via an external authentication service, and profiles are linked to Ethereum wallets.
# Stack
Golang: The main implementation language.
JWT: Used for user authentication and token validation.
MongoDB: Data storage for profile information.
Swagger: API documentation generation.
# API
1. Get Profile by ID
GET /api/profiles/{id}

Retrieves a profile based on its ID.  
Parameters:

id (path): The unique identifier of the profile.  
Response:

200: Profile data as JSON.  
400: Error message.


2. Save Profile

POST /api/profiles/
Creates a new profile.
Parameters:

profile (body): JSON object with profile information (login, name, email, etc.).   
Response:

200: Success message with access and refresh tokens.  
400: Validation error or profile already exists.   
500: Internal server error if the profile could not be saved.

3. Get Current User Profile

GET /api/profiles/me   
Retrieves the profile of the currently authenticated user.   
Response:

200: The profile of the current user.   
400: Error if the profile does not exist.

4. Update Profile

PATCH /api/profiles/   
Updates an existing profile.  
Parameters:  

profile (body): JSON object with updated profile data (login, name, email, description, photo, photoTitle).
Response:

200: Updated profile data.  
400: Validation error or login already exists.  
500: Internal server error.  

5. Get Profile by Wallet Address

POST /api/profiles/wallet-address  
Retrieves a profile based on the wallet address.  
Parameters:  

wallet_address (body): JSON object with the wallet address.  
Response:

200: The profile associated with the given wallet address.  
400: Error message.

# Swagger
Access to swagger: http://localhost:port/swagger/
# Quickstart
1. git clone <repository_url>
2. Configure the config.yml
3. docker-compose up --build 
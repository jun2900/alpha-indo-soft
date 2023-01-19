# alpha-indo-soft


Steps to install this API on local computer:

1. Download and install **docker**:
   - [Download link](https://docs.docker.com/docker-for-windows/install/) for **windows**
   - [Download link](https://docs.docker.com/docker-for-mac/install/) for **mac**
   - For linux u need to download and install docker according distribution and [docker compose](https://docs.docker.com/compose/install/)
2. Clone this repo
3. (Optional) Kill all service that are running on port 7122, 7123, and 7120 on your local computer. If there is no service that running on the specified port then skip this step
4. Build all the docker image using this command on the project root folder in your terminal:
   > docker-compose up --build
5. Now you can access the api. Here is the Postman collection regarding the api https://www.postman.com/restless-robot-876555/workspace/not-so-personal/collection/9354582-168f5b1d-d3ee-4d18-94e9-e727cec5d724?action=share&creator=9354582

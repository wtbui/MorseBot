# MorseBot

Morsebot is a Discord bot that allows users to control IoT Govee lights using commands on Discord messaging channels. MorseBot provides a seamless way to manage both individual devices and groups for synchronized lighting experiences.

## Features

- **Individual Light Control:**  
  Adjust brightness, change colors, toggle lights, and apply special effects (cycling through multiple colors/brightness settings).

- **Group Light Syncing:**  
  Control lights for groups of users on different home networks to simultaneously to enable synchronized lighting effects across multiple devices.

- **Command-Based Interaction:**  
  Users can control the lights through commands directly inputted into the messaging channels where the bot is active.

## Technical Overview

MorseBot is written in Go (Golang) and hosted on an Amazon EC2 server. It is structured into multiple modules to keep the codebase modular and maintainable.

### Project Structure

- **cmd/**  
  - Contains the entrypoint/startup code to initialize the bot using CLI.
  - Manages user Govee keys in a local database (register/delete operations).

- **pkg/**  
  - **pkg/data:**  
    Provides an API to register, delete, and parse users in the Govee key database.
  
  - **pkg/events:**  
    Acts as the main bot code interfacing with the Discord API. This module handles event processing (e.g., message events) and is the central place for adding new bot commands.
  
  - **pkg/goveego:**  
    A wrapper around the Govee API, enabling HTTP requests to control user devices.
  
  - **pkg/lightsync:**  
    Implements the "lights" command by parsing user input from text channels into executable jobs via the Govee API wrapper. This is also where new lights-based functionality can be integrated.
  
  - **pkg/utils:**  
    Contains helper functions to enhance bot functionality, including a wrapper for the Discord API's embed message function for generating user friendly bot responses.

## Future Plans/Roadmap

This is still currently being maintained with plans for future updates/features to improve scalability and performance.
### Planned Additions

- **Database:**
	Currently uses a local text based datastore. Plans to move to a relational SQL 	based database to improve lookup functionality and scalability.

- **Lights Job Execution:**
	Currently executes "lights" command jobs on a single thread, will move to a multi-threaded execution structure to improve job execution latency.


var Users: &User[];

interface User {
	Id: string;
	Name: string;
	Email: string;
}

const NewUser = {
	"Id":	 "0",
	"Name":  "New User",
	"Email": "user-new@example.com"
}

function addUser(user: User, newUser: boolean = false) {
	const userDiv = document.createElement("div");
	const nameInput = document.createElement("input");
	const emailInput = document.createElement("input");
	const deleteButton = document.createElement("button");
	const separator = document.createElement("hr")
	userDiv.id = "userDiv"
	nameInput.id = "nameInput"
	emailInput.id = "emailInput"
	deleteButton.id = "deleteButton"
	separator.id = "separator"
	var nameString = `${user.Name}`
	var emailString = `${user.Email}`
	const htmlUserList = document.getElementById("userList") as HTMLUListElement;

	deleteButton.textContent = "X";
	deleteButton.addEventListener("click", async () => {
		const userToRemove = user.Email;
		Users = Users.filter((userToDelete: User) => userToDelete.Email !== userToRemove);
		htmlUserList.removeChild(userDiv);
		htmlUserList.removeChild(separator);
		deleteData(user);
	});

	nameInput.type = "text"
	nameInput.value = nameString;
	nameInput.addEventListener("keypress", async (event: KeyboardEvent) => {
		if (event.key === "Enter") {
			event.preventDefault();
			user.Name = nameInput.value
			putData(user);
		};
	});

	emailInput.type = "text"
	emailInput.value = emailString;
	emailInput.addEventListener("keypress", async (event: KeyboardEvent) => {
		if (event.key === "Enter") {
			event.preventDefault();
			user.Email = emailInput.value
			putData(user);
		};
	});

	userDiv.appendChild(nameInput);
	userDiv.appendChild(emailInput);
	userDiv.appendChild(deleteButton);
	htmlUserList.appendChild(userDiv);
	htmlUserList.appendChild(separator);

	if (newUser) {
		postData(user);
	}
}

// Send a POST request to the Go backend to create data
async function postData(user: User) {
	try {
		const request = {
			method: "POST",
			headers: {
				"Content-Type": "application/json",
			},
			body: JSON.stringify(user),
		};

		console.log("Sending request to POST /PostData with body:", request.body);
		const response = await fetch("http://localhost:8080/PostData", request);
		const responseBody = await response.text();

		if (!response.ok) {
			throw new Error(`postData HTTP error! Status: ${response.status}`);
		}

		// Check if the response body is not empty
		if (responseBody.trim() === "") {
			console.warn("Response body is empty");
			return;
		}

		// Try to parse the JSON response
		try {
			const data = JSON.parse(responseBody);
			console.log("Response:", data);
		} catch (parseError) {
			console.error("Error parsing response body as JSON:", parseError);
			console.error("Raw response body was:", responseBody);
		}

	} catch (error) {
		console.error("Error calling postData:", error);
	}
}

// Send a GET request to the Go backend to read data
async function getData() {
	try {
		console.log("Sending request to GET /GetData");
		const response = await fetch("http://localhost:8080/GetData");

		// Throw an error if the response is not ok
		if (!response.ok) {
			throw new Error(`getData HTTP error! Status: ${response.status}`);
		}

		// Parse the JSON data from the response
		const users = await response.json();

		// Get the HTML 'userList' element
		const htmlUserList = document.getElementById("userList") as HTMLUListElement;

		if (Array.isArray(users)) {
			Users = users;
			htmlUserList.innerHTML = "";
			users.forEach((user: { Id: string, Name: string, Email: string }) => {
				addUser(user);
			});
			console.log("getData request successful: ", users)
		}
	} catch (error) {
		console.error("Error fetching user data:", error);
	}
}

// Send a PUT request to the Go backend to update data
async function putData(user: User) {
	try {
		const request = {
			method: "PUT",
			headers: {
				"Content-Type": "application/json",
			},
			body: JSON.stringify(user),
		};

		console.log("Sending request to PUT /PutData with body:", request.body);
		const response = await fetch("http://localhost:8080/PutData", request);
		const responseBody = await response.text();
		console.log("Response body:", responseBody);

		// Throw an error if the response is not ok
		if (!response.ok) {
			throw new Error(`putData HTTP error! Status: ${response.status}`);
		}

		// Parse the JSON data from the response
		const data = JSON.parse(responseBody);
		console.log("Response from Express API:", data);
	} catch (error) {
		console.error("Error calling putData:", error);
	}
}

// Send a DELETE request to the go backend to delete data
async function deleteData(user: User) {
	try {
		const request = {
			method: "DELETE",
			headers: {
				"Content-Type": "application/json",
			},
			body: JSON.stringify(user)
		};

		const response = await fetch("http://localhost:8080/DeleteData", request);

		// Throw an error if the response is not ok
		if (!response.ok) {
			throw new Error(`deleteData HTTP error! Status: ${response.status}`);
		}

		let data = {};
		const responseBody = await response.text();

		if (responseBody) {
			data = JSON.parse(responseBody);
		}

		console.log("Response:", data);
	} catch (error) {
		console.error("Error calling deleteData:", error);
	}
}

window.addEventListener("DOMContentLoaded", () => {
	getData();

	const loadDataButton = document.getElementById("loadDataButton") as HTMLButtonElement;
	const addUserButton = document.getElementById("addUserButton") as HTMLButtonElement;

	addUserButton.addEventListener("click", (event: MouseEvent) => {
		event.preventDefault();
		var newUser = NewUser;
		addUser(newUser, true);
	});

	loadDataButton.addEventListener("click", async (event: MouseEvent) => {
		event.preventDefault();

		try {
			await getData();
		} catch (error) {
			console.error(error);
		}
	});
});
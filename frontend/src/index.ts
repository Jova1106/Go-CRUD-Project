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

async function postData(user: User) {
	try {
		const request = {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
			},
			body: JSON.stringify(user)
		};

		const response = await fetch('http://localhost:8080/PostData', request);

		if (!response.ok) {
			throw new Error(`HTTP error! Status: ${response.status}`);
		}

		const data = await response.json();
		console.log('Response:', data);
	} catch (error) {
		console.error('Error calling postData:', error);
	}
}

async function putData(user: User) {
	try {
		const request = {
			method: 'PUT',
			headers: {
				'Content-Type': 'application/json',
			},
			body: JSON.stringify(user),
		};

		console.log('Sending request to PUT /PutData with body:', request.body);
		const response = await fetch('http://localhost:8080/PutData', request);
		const responseBody = await response.text();
		console.log('Response body:', responseBody);

		if (!response.ok) {
			console.error(`HTTP error! Status: ${response.status}`);
			throw new Error(`HTTP error! Status: ${response.status}`);
		}

		const data = responseBody ? JSON.parse(responseBody) : null;
		console.log('Response from Express API:', data);
	} catch (error) {
		console.error('Error calling putData:', error);
	}
}

async function deleteData(user: User) {
	try {
		const request = {
			method: 'DELETE',
			headers: {
				'Content-Type': 'application/json',
			},
			body: JSON.stringify(user)
		};

		const response = await fetch('http://localhost:8080/DeleteData', request);

		if (!response.ok) {
			throw new Error(`HTTP error! Status: ${response.status}`);
		}

		let data = {};
		const text = await response.text();

		if (text) {
			data = JSON.parse(text);
		}

		console.log('Response:', data);
	} catch (error) {
		console.error('Error calling deleteData:', error);
	}
}


function addUser(user: User, newUser: boolean = false) {
	const userDiv = document.createElement('div');
	const nameInput = document.createElement('input');
	const emailInput = document.createElement('input');
	const deleteButton = document.createElement('button');
	const separator = document.createElement('hr')
	var nameString = `${user.Name}`
	var emailString = `${user.Email}`
	const htmlUserList = document.getElementById('userList') as HTMLUListElement;

	deleteButton.textContent = "X";
	deleteButton.addEventListener('click', async () => {
		const userToRemove = user.Email;
		Users = Users.filter((userToDelete: User) => userToDelete.Email !== userToRemove);
		htmlUserList.removeChild(userDiv);
		htmlUserList.removeChild(separator);
		deleteData(user);
	});

	Object.assign(userDiv.style, {
		background: "#b8bfc5",
		margin: "5px 0",
		padding: "10px",
		borderRadius: "5px",
		display: "flex",
	});

	Object.assign(nameInput.style, {
		background: "#e9ecef",
		float: "left",
		margin: "5px 0",
		padding: "10px",
		borderRadius: "5px",
		display: "flex",
		flexGrow: 1
	});

	Object.assign(emailInput.style, {
		background: "#e9ecef",
		float: "left",
		margin: "5px 0",
		padding: "10px",
		borderRadius: "5px",
		display: "flex",
		flexGrow: 1
	});

	Object.assign(deleteButton.style, {
		marginLeft: "10px",
		padding: "5px 10px",
		borderRadius: "5px",
		display: "flex"
	});

	nameInput.type = 'text'
	nameInput.value = nameString;
	nameInput.addEventListener('keypress', async (event: KeyboardEvent) => {
		if (event.key === 'Enter') {
			event.preventDefault();
			user.Name = nameInput.value
			putData(user);
		};
	});

	emailInput.type = 'text'
	emailInput.value = emailString;
	emailInput.addEventListener('keypress', async (event: KeyboardEvent) => {
		if (event.key === 'Enter') {
			event.preventDefault();
			user.Email = emailInput.value
			putData(user);
		};
	});

	userDiv.appendChild(nameInput);
	userDiv.appendChild(emailInput);
	userDiv.appendChild(deleteButton);
	htmlUserList.appendChild(userDiv);
	htmlUserList.appendChild(separator);3
	if (newUser) {
		postData(user);
	}
}

function loadData() {
	fetch('http://localhost:8080/GetData')
	.then(response => response.json())
	.then(users => {
		Users = users
		const htmlUserList = document.getElementById('userList') as HTMLUListElement;
		if (Array.isArray(Users)) {
			htmlUserList.innerHTML = '';
			users.forEach((user: {Id: string, Name: string, Email: string}) => {
				addUser(user)
			})
		}
	})
	.catch(error => {
		console.error('Error fetching user data:', error);
	});
}

window.addEventListener('DOMContentLoaded', () => {
	loadData();

	const loadDataButton = document.getElementById('loadDataButton') as HTMLButtonElement;
	const addUserButton = document.getElementById('addUserButton') as HTMLButtonElement;

	addUserButton.addEventListener('click', (event: MouseEvent) => {
		event.preventDefault();
		var newUser = NewUser;
		addUser(newUser, true);
	});

	loadDataButton.addEventListener('click', async (event: MouseEvent) => {
		event.preventDefault();

		const promises = Users.map(user => putData(user));

		try {
			await Promise.all(promises);
			await loadData();
		} catch (error) {
			console.error("Error saving data before refresh:", error);
		}
	});
});
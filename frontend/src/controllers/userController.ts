import {Request, Response} from 'express';
import axios from 'axios';

const GO_BACKEND_URL = 'http://localhost:8080/';

export const postData = async (req: Request, res: Response): Promise<void> => {
	try {
		const { userToPost } = req.body;

		if (!userToPost) {
			res.status(400).json({ message: 'userToPost is required' });
			return;
		}

		const response = await axios.post(`${GO_BACKEND_URL}PostData`, {
			params: {userToPost}
		});

		res.json(response.data);
	} catch (error) {
		console.error('Error updating data:', error);
		res.status(500).json({ message: 'Error updating data' });
	}
};

export const putData = async (req: Request, res: Response): Promise<void> => {
	try {
		const { userToPut } = req.body;

		if (!userToPut) {
			res.status(400).json({ message: 'userToPut is required' });
			return;
		}

		const response = await axios.put(`${GO_BACKEND_URL}PutData`, {
			params: {userToPut}
		});

		res.json(response.data);
	} catch (error) {
		console.error('Error updating data:', error);
		res.status(500).json({ message: 'Error updating data' });
	}
};

export const getData = async (req: Request, res: Response) => {
	try {
		const response = await axios.get(GO_BACKEND_URL + "GetData");
		res.json(response.data);
	} catch (error) {
		console.error('Error fetching data from Go API:', error);
		res.status(500).json({ message: 'Error fetching data from Go API' });
	}
};

export const deleteData = async (req: Request, res: Response): Promise<void> => {
	try {
		const { userToDelete } = req.body;

		if (!userToDelete) {
			res.status(400).json({ message: 'userToDelete is required' });
			return;
		}

		const response = await axios.delete(`${GO_BACKEND_URL}DeleteData`, {
			params: {userToDelete}
		});

		res.json(response.data);
	} catch (error) {
		console.error('Error updating data:', error);
		res.status(500).json({ message: 'Error updating data' });
	}
};
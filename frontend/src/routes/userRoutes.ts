// src/routes/userRoutes.ts
import express from 'express';
import {postData, putData, getData, deleteData} from '../controllers/userController';

const router = express.Router();

router.post('/PostData', postData);
router.put('/PutData', putData);
router.get('/GetData', getData);
router.delete('/DeleteData', deleteData);

export default router;
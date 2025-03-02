import express, {Request, Response} from 'express';
import userRoutes from './routes/userRoutes';
import path from 'path';

const app = express();
app.use(express.static(path.join(__dirname, 'views')));
app.use(express.json());
app.use('/', userRoutes);

app.get('/', (req: Request, res: Response) => {
	res.sendFile(path.join(__dirname, 'views', '../../index.html'))
})

export default app;

import dotenv from 'dotenv';
import express, { Express, Request, Response } from 'express';
import { AuthMiddleware } from './middleware/auth.middleware';
import { UserService } from './services/user.service';

dotenv.config();

const app: Express = express();
const port = process.env.PORT || 3000;

app.use(express.json());
app.use(AuthMiddleware);

const userService = new UserService();

app.get('/api/users', async (req: Request, res: Response) => {
    try {
        const users = await userService.getAllUsers();
        res.json(users);
    } catch (error) {
        res.status(500).json({ error: 'Failed to fetch users' });
    }
});

app.post('/api/users', async (req: Request, res: Response) => {
    try {
        const user = await userService.createUser(req.body);
        res.status(201).json(user);
    } catch (error) {
        res.status(400).json({ error: 'Failed to create user' });
    }
});

app.listen(port, () => {
    console.log(`Server running at http://localhost:${port}`);
});

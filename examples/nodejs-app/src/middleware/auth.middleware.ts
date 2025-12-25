import { NextFunction, Request, Response } from 'express';

export interface AuthRequest extends Request {
    userId?: string;
    token?: string;
}

export const AuthMiddleware = (req: AuthRequest, res: Response, next: NextFunction) => {
    const token = req.headers.authorization?.split(' ')[1];

    if (!token) {
        return res.status(401).json({ error: 'No token provided' });
    }

    try {
        // Simplified token validation
        const decoded = Buffer.from(token, 'base64').toString('utf-8');
        req.userId = decoded;
        req.token = token;
        next();
    } catch (error) {
        res.status(401).json({ error: 'Invalid token' });
    }
};

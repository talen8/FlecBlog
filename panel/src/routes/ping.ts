import { Hono } from 'hono';
import { Env } from '../types';

const pingApi = new Hono<{ Bindings: Env }>();

pingApi.get('/', async (c) => c.json({ status: 'ok' }));

export { pingApi };

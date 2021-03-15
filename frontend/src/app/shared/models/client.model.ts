import { Job } from './job.model'

export class Client {
    public id: number;
    public name: string
    public phone: string;
    public email: string;
    public jobs: Job[];
  }
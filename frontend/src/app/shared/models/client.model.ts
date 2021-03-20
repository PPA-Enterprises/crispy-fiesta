import { Job } from './job.model'
import { Log } from './log.model'


export class Client {
    public _id: string;
    public name: string
    public phone: string;
    public email: string;
    public jobs: Job[];
    public history: Log[];
  }
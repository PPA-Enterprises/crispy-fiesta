import { Job } from './job.model'

export class Log {
    public id: string;
    public event_type: string;
    public timestamp: Date;
    public editor: string;
    public editor_id: string;
    public persisted: Boolean;
    public changes: any;
  }
export class Job {
    public _id?: string;
    public client_name: string
    public client_phone: string;
    public car_info: string;
    public notes: string;
    public tag: string;
    public start: Date;
    public end: Date;
    public title: string;
    public primary_color: string;
    public secondary_color: string;
    public history?: any;
  }
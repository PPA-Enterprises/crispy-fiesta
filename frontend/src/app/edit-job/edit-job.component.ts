import { Component, ViewChild, OnInit } from '@angular/core';
import { FormControl, FormGroup, Validators, NgForm } from '@angular/forms';
import { Router, ActivatedRoute, ParamMap } from '@angular/router';
import { JobService } from '../shared/services/job.service'
import { Job } from '../shared/models/job.model';
import { startOfDay, subDays, addDays} from 'date-fns';

export class JobForm {
  public fname: string;
  public lname: string;
  public phone: string;
  public carInfo: string;
  public notes: string;
  public primary: any;
  public secondary: any;
}

@Component({
  selector: 'app-page',
  templateUrl: './edit-job.component.html',
  styleUrls: ['./edit-job.component.scss']
})

export class EditJobComponent implements OnInit {
  public job: Job;
  id: string;
  private sub: any;
  model = new JobForm();
  submitted = false;
  name: string[];
  notClosed = true;
  startDate: any;
  endDate: any;


  constructor(private jobService: JobService, private route: ActivatedRoute, private router: Router) {}

  ngOnInit(): void {
    this.sub = this.route.params.subscribe(params => {
      this.id = params['id'];
    });

    this.jobService.getJobById(this.id).subscribe((incomingJob: Job) => {
      this.job = incomingJob;

      if(this.job.tag == "CLOSED") {
        this.notClosed = false;
      }

      this.name = this.job.client_name.split(" ");
      this.model = {
        fname: this.name[0],
        lname: this.name[1],
        phone: this.job.client_phone,
        carInfo: this.job.car_info,
        notes: this.job.notes,
        primary: this.job.color.primary,
        secondary: this.job.color.secondary
      }

      this.startDate = this.job.start;
      this.endDate = this.job.end;
    })
  }

  ngOnDestroy() {
    this.sub.unsubscribe();
  }

  onSubmit(form) {
    //console.log(this.startDate);
    
    this.job = {
      _id: this.job._id,
      client_name: this.model.fname + " " + this.model.lname,
      client_phone: this.model.phone,
      car_info: this.model.carInfo,
      notes: this.model.notes,
      tag: "OPEN",
      start: this.startDate,
      end: this.endDate,
      title: this.model.carInfo + " - " + this.model.fname + " " + this.model.lname,
      color: {
        primary: this.model.primary,
        secondary: this.model.secondary,
      }
    }

    this.jobService.editJobById(this.id, this.job).subscribe((job: Job) => {
      this.job = job
      this.router.navigate(['/jobs']);
    })
  }

  sendToTintWork() {
    console.log("Before:" + this.job);
    this.job = {
      _id: this.job._id,
      client_name: this.model.fname + " " + this.model.lname,
      client_phone: this.model.phone,
      car_info: this.model.carInfo,
      notes: this.model.notes,
      tag: "CLOSED",
      start: subDays(startOfDay(new Date()), 1),
      end: addDays(new Date(), 1),
      title: this.model.carInfo + " - " + this.model.fname + " " + this.model.lname,
      color: {
        primary: '#ad2121',
        secondary: '#FAE3E3',
      }
    }

    this.jobService.editJobById(this.id, this.job).subscribe((job: Job) => {
      this.job = job
      this.router.navigate(['/jobs']);
    })

  }

  delete() {
    this.jobService.deleteJobById(this.id).subscribe((deleted: any) => {
      if (deleted) this.router.navigate(['/jobs']);
    })
  }

}

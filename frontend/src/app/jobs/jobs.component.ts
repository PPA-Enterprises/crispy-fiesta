import { Component, ViewChild, OnInit } from '@angular/core';
import { FormControl, FormGroup, Validators, NgForm } from '@angular/forms';
import { JobService } from '../shared/services/job.service'

// export class JobForm {
//   public fname: string;
//   public lname: string;
//   public phone: string;
//   public carInfo: string;
//   public apptInfo: string;
//   public notes: string;
// }

// export class Job {
//   public client_name: string
//   public client_phone: string;
//   public car_info: string;
//   public appointment_info: string;
//   public notes: string;
// }

@Component({
  selector: 'app-page',
  templateUrl: './jobs.component.html',
  styleUrls: ['./jobs.component.scss']
})

export class JobsComponent implements OnInit {
  newJobs: any[] = [
    { name: "Tristan Hull", date: "December 12th, 2020", car: "2003 Saturn Vue" }, { name: "Joshua Benz", date: "December 17th, 2020", car: "2010 Ford Focus" }, { name: "Frank Swartz", date: "December 6th, 2020", car: "2019 Ford Mustang" },
  ];

  constructor(private jobService: JobService) {
   }

  ngOnInit() {

  }

}

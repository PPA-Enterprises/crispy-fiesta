import { Component, ViewChild, OnInit } from '@angular/core';
import { FormControl, FormGroup, Validators, NgForm } from '@angular/forms';
import { Router, ActivatedRoute, ParamMap } from '@angular/router';
import { JobService } from '../shared/services/job.service'
import { ClientService } from '../shared/services/client.service'
import { Job } from '../shared/models/job.model'
import { Client } from '../shared/models/client.model'


export class ClientForm {
  public lname: string
  public fname: string
  public phone: string;
  public email: string;
}

@Component({
  selector: 'app-page',
  templateUrl: './client.component.html',
  styleUrls: ['./client.component.scss']
})

export class ClientComponent implements OnInit {
  id: string;
  private sub: any;
  client: Client;
  model = new ClientForm();
  name: string[];

  constructor(private jobService: JobService, private clientService: ClientService, private route: ActivatedRoute, private router: Router) {}

  ngOnInit() {
    this.sub = this.route.params.subscribe(params => {
      this.id = params['id'];
    });

    this.clientService.getClientById(this.id).subscribe((client: Client) => {
    this.client = client;
    this.name = this.client.name.split(" ");
    this.model = {
      fname: this.name[0],
      lname: this.name[1],
      phone: this.client.phone,
      email: this.client.email
    }
    })
  }

  delete() {

  }

  onSubmit(form) {
    this.client.name = this.model.fname + " " + this.model.lname;
    this.client.email = this.model.email;
    this.client.phone = this.model.phone;

    this.clientService.editClientById(this.id, this.client).subscribe((client: Client) => {
      this.client = client;
      this.router.navigate(['/clients']);
    });
  }


}

import { NgModule } from '@angular/core';
import { CommonModule } from "@angular/common";

import { CreateJobRoutingModule } from "./create-job-routing.module";

import { CreateJobComponent } from "./create-job.component";


@NgModule({
  imports: [
    CommonModule,
    CreateJobRoutingModule
  ],
  exports: [],
  declarations: [
    CreateJobComponent
  ],
  providers: [],
})
export class CreateJobModule { }

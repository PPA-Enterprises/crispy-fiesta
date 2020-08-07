import { NgModule } from '@angular/core';
import { CommonModule } from "@angular/common";
import { FormsModule, ReactiveFormsModule } from '@angular/forms';

import { JobsRoutingModule } from "./jobs-routing.module";

import { JobsComponent } from "./jobs.component";


@NgModule({
  imports: [
    CommonModule,
    JobsRoutingModule,
    ReactiveFormsModule,
    FormsModule,
  ],
  exports: [],
  declarations: [
    JobsComponent,
  ],
  providers: [],
})
export class JobsModule { }

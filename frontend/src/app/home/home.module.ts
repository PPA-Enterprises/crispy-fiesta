import { NgModule } from '@angular/core';
import { CommonModule } from "@angular/common";

import { HomeRoutingModule } from "./home-routing.module";

import { HomeComponent } from "./home.component";


@NgModule({
  imports: [
    CommonModule,
    HomeRoutingModule
  ],
  exports: [],
  declarations: [
    HomeComponent
  ],
  providers: [],
})
export class HomeModule { }

import { Component, ChangeDetectionStrategy, ViewChild, TemplateRef, ViewEncapsulation } from '@angular/core';
import { Router } from '@angular/router';
import { startOfDay, endOfDay, subDays, addDays, endOfMonth, isSameDay, isSameMonth, addHours } from 'date-fns';
import { Subject } from 'rxjs';
import { NgbModal } from '@ng-bootstrap/ng-bootstrap';
import { CalendarEvent, CalendarEventAction, CalendarEventTimesChangedEvent, CalendarView, } from 'angular-calendar';
import { JobService } from '../shared/services/job.service'
import { ClientService } from '../shared/services/client.service'
import { Job } from '../shared/models/job.model'
import { Client } from '../shared/models/client.model'


@Component({
  selector: 'app-page',
  changeDetection: ChangeDetectionStrategy.OnPush,
  templateUrl: './cal.component.html',
  styleUrls: ['./cal.component.scss'],
  //encapsulation: ViewEncapsulation.None
})

export class CalComponent {
  @ViewChild('modalContent', { static: true }) modalContent: TemplateRef<any>;

  constructor(private modal: NgbModal, private jobService: JobService, private router: Router) {}

  view: CalendarView = CalendarView.Month;
  CalendarView = CalendarView;
  viewDate: Date = new Date();
  modalData: {
    action: string;
    event: CalendarEvent;
  };
  refresh: Subject<any> = new Subject();
  events: CalendarEvent[] = this.jobService.getAllJobs();
  activeDayIsOpen: boolean = true;



  dayClicked({ date, events }: { date: Date; events: CalendarEvent[] }): void {
    if (isSameMonth(date, this.viewDate)) {
      if (
        (isSameDay(this.viewDate, date) && this.activeDayIsOpen === true) ||
        events.length === 0
      ) {
        this.activeDayIsOpen = false;
      } else {
        this.activeDayIsOpen = true;
      }
      this.viewDate = date;
    }
  }


  handleEvent(action: string, event: CalendarEvent): void {
    this.modalData = { event, action };
    this.modal.open(this.modalContent, { size: 'lg' });

    //this.router.navigateByUrl('/edit-job/'+event.id);
  }


  setView(view: CalendarView) {
    this.view = view;
  }

  closeOpenMonthViewDay() {
    this.activeDayIsOpen = false;
  }

}

package main;

import (
	"net/http"
	"log"
	"encoding/json"
)

type event struct {
	Name, Email, Location, Date, Time, Title, Description, Image string
}

type update struct {
	Old event
	New event
}

func handleEvent(w http.ResponseWriter, req *http.Request) {
	var user User;
	var e event;
	if req.Method != http.MethodGet && req.Method != http.MethodPut {
		json.NewDecoder(req.Body).Decode(&e);
		db.Where(&User{ Email: e.Email }).First(&user);
		if req.Method == http.MethodPost {
			user.postEvent(e);
			successRequest(w, "added event", "added event");
		} else if req.Method == http.MethodDelete {
			user.deleteEvent(e);
			successRequest(w, "removed event successfully", "deleted event");
		} else if req.Method != http.MethodOptions {
			log.Println("request method not supported");
		}
	} else if req.Method == http.MethodGet {
			email := req.Header.Get("Email");
			if len(email) > 0 {
				db.Where(&User{Email: email}).First(&user);
				events := user.getEvents(e);
				r, _ := json.Marshal(events);
				w.Write(r);
			} else {
				badRequest(w, "email not found", 400)
			}
	} else {
		var u update;
		json.NewDecoder(req.Body).Decode(&u);
		db.Where(&User{ Email: u.Old.Email }).First(&user);
		user.updateEvent(u);
		successRequest(w, "updated event", "sucessfully updated event");
	}
}

func (u User) postEvent(e event) {
	db.Create(&Event{UserID: u.ID, Location: e.Location, Date: e.Date, Time: e.Time, Title: e.Title, Description: e.Description, Image: e.Image});
}

func (u User) deleteEvent(e event) {
	var ev Event;
	db.Where(&Event{UserID: u.ID, Location: e.Location, Date: e.Date, Time: e.Time, Title: e.Title, Description: e.Description, Image: e.Image }).First(&ev);
	if ev.UserID > 0 {
		db.Delete(&ev);
	}
}

func (u User) getEvents(e event) []event {
	var events []event;
	if u.ID > 0 {
		friends := u.findFriends();
		for _, f := range friends {
			var evt []Event;
			var res []event;
			var user User;
			db.Where(&User{Email: f.Email}).First(&user);
			if user.ID > 0 {
				db.Where(&Event{UserID: user.ID}).Find(&evt);
				for _, e := range evt {
					ev := event{
						f.Name,
						f.Email,
						e.Location,
						e.Date,
						e.Time,
						e.Title,
						e.Description,
						e.Image,
					}
					res = append(res, ev);
				}
				events = append(events, res...);
			}
		}
	} 
	return events;
}

func (user User) updateEvent(u update) {
	var ev Event;
	db.Where(&Event{UserID: user.ID, Location: u.Old.Location, Date: u.Old.Date, Time: u.Old.Time, Title: u.Old.Title, Description: u.Old.Description, Image: u.Old.Image }).First(&ev);
	ev.Title = u.New.Title;
	ev.Date = u.New.Date;
	ev.Location = u.New.Location;
	ev.Time = u.New.Time;
	ev.Description = u.New.Description;
	ev.Image = u.New.Image;
	db.Save(&ev);
}
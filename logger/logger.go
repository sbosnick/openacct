// Copyright Steven Bosnick 2016. All rights reserved.
// Use of this source code is governed by the GNU General Public License version 3.
// See the file COPYING for your rights under that license.

// Package logger provides the logging configurations for other packages in the project.
// Each packages should declarae a variable named 'log' which it then uses for logging.
// For package 'mypackage' this should be declared as follows:
//	var (
//	    log = logger.New("mypackage")
//	)
// Currently the logger returned is a default logger in all instances, but this will
// change to allow a more configurable logging.
package logger

import (
	"github.com/Sirupsen/logrus"
	"github.com/gogap/logrus_mate"
)

// Returns a logger for the package 'packageName'
func New(packageName string) *logrus.Logger {
	return logrus_mate.Logger()
}
